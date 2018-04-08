package kubernetestracker

import (
	"errors"
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/helper"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type KubernetesTracker struct {
	clientSet *kubernetes.Clientset
}

func New(cs *kubernetes.Clientset) (*KubernetesTracker, error) {
	if cs == nil {
		var err error
		cs, err = NewClientSet()
		if err != nil {
			return nil, err
		}
	}
	return &KubernetesTracker{
		clientSet: cs,
	}, nil
}

func (kt *KubernetesTracker) ListJobCategories() ([]string, error) {
	return []string{}, nil
}

func (kt *KubernetesTracker) ListJobs() ([]string, error) {
	jc, err := getJobsClient(kt.clientSet)
	if err != nil {
		return nil, fmt.Errorf("ListJobs: %s", err.Error())
	}
	jobsList, err := jc.List(k8sapi.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing jobs with client: %s", err.Error())
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
		return "", fmt.Errorf("converting job template into a k8s job: %s", err.Error())
	}
	jc, err := getJobsClient(kt.clientSet)
	if err != nil {
		return "", fmt.Errorf("get client: %s", err.Error())
	}
	j, err := jc.Create(job)
	if err != nil {
		return "", fmt.Errorf("creating new job: %s", err.Error())
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
	jc, err := getJobsClient(kt.clientSet)
	if err != nil {
		return drmaa2interface.Undetermined
	}
	return DRMAA2State(jc, jobid)
}

func (kt *KubernetesTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	jc, err := getJobsClient(kt.clientSet)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}
	return JobToJobInfo(jc, jobid)
}

func (kt *KubernetesTracker) JobControl(jobid, state string) error {
	jc, err := getJobsClient(kt.clientSet)
	if err != nil {
		return fmt.Errorf("JobControl: %s", err.Error())
	}

	job, err := getJobByID(jc, jobid)
	if err != nil {
		return err
	}

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
		return jc.Delete(job.GetName(), &k8sapi.DeleteOptions{})
	}
	return errors.New("undefined state")
}

func (kt *KubernetesTracker) Wait(jobid string, timeout time.Duration, states ...drmaa2interface.JobState) error {
	return helper.WaitForState(kt, jobid, timeout, states...)
}

func (kt *KubernetesTracker) DeleteJob(jobid string) error {
	return nil
}
