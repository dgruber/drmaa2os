package kubernetestracker

import (
	"errors"
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/helper"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type KubernetesTracker struct{}

func New() (*KubernetesTracker, error) {
	return &KubernetesTracker{}, nil
}

func (kt *KubernetesTracker) ListJobCategories() ([]string, error) {
	// external registry
	return nil, nil
}

func (kt *KubernetesTracker) ListJobs() ([]string, error) {
	cs, err := CreateClientSet()
	if err != nil {
		return nil, fmt.Errorf("error during addjob client creation: %s", err.Error())
	}
	jobsClient := cs.BatchV1().Jobs("default")
	jobsList, err := jobsClient.List(k8sapi.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error during listing jobs with client: %s", err.Error())
	}
	ids := make([]string, 0, len(jobsList.Items))
	for _, job := range jobsList.Items {
		ids = append(ids, string(job.UID))
	}
	return ids, nil
}

func (kt *KubernetesTracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	job, err := convertJob(jt)
	if err != nil {
		return "", fmt.Errorf("error during converting job template into a k8s job: %s", err.Error())
	}
	cs, err := CreateClientSet()
	if err != nil {
		return "", fmt.Errorf("error during addjob client creation: %s", err.Error())
	}
	jc := cs.BatchV1().Jobs("default")
	j, err := jc.Create(job)
	if err != nil {
		return "", fmt.Errorf("error during k8s job client initialization: %s", err.Error())
	}
	return string(j.UID), nil
}

func (kt *KubernetesTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return helper.AddArrayJobAsSingleJobs(jt, kt, begin, end, step)
}

func (kt *KubernetesTracker) ListArrayJobs(id string) ([]string, error) {
	return helper.ArrayJobID2GUIDs(id)
}

func (kt *KubernetesTracker) JobState(jobid string) drmaa2interface.JobState {
	cs, err := CreateClientSet()
	if err != nil {
		return drmaa2interface.Undetermined
	}
	jc := cs.BatchV1().Jobs("default")

	job, err := jc.Get(jobid, meta_v1.GetOptions{})
	if err != nil || job == nil {
		return drmaa2interface.Undetermined
	}
	return convertJobStatus2JobState(&job.Status)
}

func (kt *KubernetesTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	return drmaa2interface.JobInfo{}, nil
}

func (kt *KubernetesTracker) JobControl(jobid, state string) error {
	cs, err := CreateClientSet()
	if err != nil {
		return err
	}
	jc := cs.BatchV1().Jobs("default")

	switch state {
	case "suspend":
		return errors.New("Unsupported Operation")
	case "resume":
		return errors.New("Unsupported Operation")
	case "hold":
		return errors.New("Unsupported Operation")
	case "release":
		return errors.New("Unsupported Operation")
	case "terminate":
		return jc.Delete(jobid, &meta_v1.DeleteOptions{})
	}
	return errors.New("undefined state")
}

func (kt *KubernetesTracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return nil
}

func (kt *KubernetesTracker) DeleteJob(jobid string) error {
	return nil
}
