package containerdtracker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/dgruber/drmaa2interface"
)

// ContainerDTracker is an implementation of the drmaa2.JobTracker interface for containerd
type ContainerDTracker struct {
	client *containerd.Client
	ns     string
}

// NewJobTracker creates a new JobTracker instance
func NewJobTracker(client *containerd.Client, namespace string) *JobTracker {
	return &ContainerDTracker{
		client: client,
		ns:     namespace,
	}
}

// JobState returns the current state of the job with the given ID
func (t *ContainerDTracker) JobState(id string) (drmaa2.JobState, error) {
	job, err := t.GetJob(id)
	if err != nil {
		return drmaa2.Undetermined, err
	}
	return job.State()
}

// JobControl allows control actions to be performed on the job with the given ID
func (t *ContainerDTracker) JobControl(id string, action drmaa2.JobControlAction) error {
	job, err := t.GetJob(id)
	if err != nil {
		return err
	}
	return job.Control(action)
}

func (dt *ContainerDTracker) ListJobs() ([]string, error) {
	return nil, nil
}

func (dt *ContainerDTracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	return "", nil
}

func (dt *ContainerDTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return "", nil
}

func (dt *ContainerDTracker) ListArrayJobs(string) ([]string, error) {
	return nil, nil
}

func (dt *ContainerDTracker) JobState(jobid string) (drmaa2interface.JobState, string, error) {
	return drmaa2interface.Undetermined, "", nil
}

func (dt *ContainerDTracker) JobInfo(jobID string) (drmaa2interface.JobInfo, error) {
	ctx := namespaces.WithNamespace(context.Background(), dt.ns)
	container, err := dt.client.LoadContainer(ctx, jobID)
	if err != nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("could not load container: %v", err)
	}
	return ConvertContainerToInfo(ctx, container)
}

func (dt *ContainerDTracker) JobControl(jobid, state string) error {
	switch state {
	case "suspend":
		return dt.cli.ContainerKill(context.Background(), jobid, "SIGSTOP")
	case "resume":
		return dt.cli.ContainerKill(context.Background(), jobid, "SIGCONT")
	case "hold":
		return errors.New("Unsupported Operation")
	case "release":
		return errors.New("Unsupported Operation")
	case "terminate":
		return dt.cli.ContainerKill(context.Background(), jobid, "SIGKILL")
	}
	return errors.New("undefined state")
	return nil
}

func (dt *ContainerDTracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return nil
}

func (dt *ContainerDTracker) DeleteJob(jobid string) error {
	return nil
}
