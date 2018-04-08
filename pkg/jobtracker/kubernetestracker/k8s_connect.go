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

func NewClientSet() (*kubernetes.Clientset, error) {
	kubeconfig, err := kubeConfigFile()
	if err != nil {
		return nil, fmt.Errorf("reading .kube/config file: %s", err.Error())
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("reading .kube/config file: %s", err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("reading .kube/config file: %s", err.Error())
	}
	return clientSet, nil
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

func getJobsClient(cs *kubernetes.Clientset) (batchv1.JobInterface, error) {
	return cs.BatchV1().Jobs("default"), nil
}
