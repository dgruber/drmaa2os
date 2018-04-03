package kubernetestracker

import (
	"errors"
	"fmt"
	"k8s.io/api/batch/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var (
	clientSet *kubernetes.Clientset
)

func init() {
	kubeconfig, err := kubeConfigFile()
	if err != nil {
		panic(err)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func kubeConfigFile() (string, error) {
	home := homeDir()
	if home == "" {
		return "", errors.New("home dir not found")
	}
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	if _, err := os.Stat(kubeconfig); err != nil {
		return "", errors.New("home does not contain .kube config file")
	}
	return kubeconfig, nil
}

func CreateClientSet() (*kubernetes.Clientset, error) {
	return clientSet, nil
}

func getJobByID(jc batchv1.JobInterface, jobid string) (*v1.Job, error) {
	jobs, err := jc.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, job := range jobs.Items {
		if jobid == string(job.GetUID()) {
			return &job, nil
		}
	}
	return nil, fmt.Errorf("job with jobid %s not found", jobid)
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	return os.Getenv("USERPROFILE")
}

func getJobsClient() (batchv1.JobInterface, error) {
	cs, err := CreateClientSet()
	if err != nil {
		return nil, fmt.Errorf("jobs client creation: %s", err.Error())
	}
	return cs.BatchV1().Jobs("default"), nil
}
