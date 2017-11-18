package kubernetestracker

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"time"
)

type KubernetesTracker struct{}

func New() (*KubernetesTracker, error) {
	return &KubernetesTracker{}, nil
}

func (kt *KubernetesTracker) ListJobs() ([]string, error) {
	return nil, nil
}

func (kt *KubernetesTracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	_, err := convertJob(jt)
	if err != nil {
		return "", fmt.Errorf("error during converting job template into a k8s job: %s", err.Error())
	}
	return "", nil
}

func (kt *KubernetesTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return "", nil
}

func (kt *KubernetesTracker) ListArrayJobs(string) ([]string, error) {
	return nil, nil
}

func (kt *KubernetesTracker) JobState(jobid string) drmaa2interface.JobState {
	return drmaa2interface.Undetermined
}

func (kt *KubernetesTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	return drmaa2interface.JobInfo{}, nil
}

func (kt *KubernetesTracker) JobControl(jobid, state string) error {
	return nil
}

func (kt *KubernetesTracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return nil
}

func (kt *KubernetesTracker) DeleteJob(jobid string) error {
	return nil
}
