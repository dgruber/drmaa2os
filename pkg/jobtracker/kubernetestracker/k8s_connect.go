package kubernetestracker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // support for GCP / GKE
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientSet create a new clientset by parsing the .kube/config file
// in the home directory.
func NewClientSet() (*kubernetes.Clientset, error) {
	// when running inside Kubernetes as pod we need to use rest
	// rather than kubeconfig
	if _, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST"); exists == true {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("expecting running inside a pod but do not get cluster config: %v", err)
		}
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("expecting running inside a pod but do not get client set: %v", err)
		}
		return clientSet, nil
	}

	kubeconfig, err := kubeConfigFile()
	if err != nil {
		return nil, fmt.Errorf("opening .kube/config file: %s", err.Error())
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("reading .kube/config file: %s", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating ClientSet from .kube/config file: %s", err.Error())
	}
	return clientSet, nil
}

func kubeConfigFile() (string, error) {
	// let KUBECONFIG override the default
	if env, exists := os.LookupEnv("KUBECONFIG"); exists && env != "" {
		return env, nil
	}
	home := homeDir()
	if home == "" {
		return "", errors.New("home dir not found")
	}
	kubeconfig := filepath.Join(home, ".kube", "config")
	if _, err := os.Stat(kubeconfig); err != nil {
		return "", errors.New("home does not contain .kube config file")
	}
	return kubeconfig, nil
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE")
}

func getJobsClient(cs *kubernetes.Clientset, namespace string) (batchv1.JobInterface, error) {
	return cs.BatchV1().Jobs(namespace), nil
}
