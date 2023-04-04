package containerdtracker

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/namespaces"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/helper"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
)

// ContainerdJobTracker implements the JobTracker interface for containerd.
type ContainerdJobTracker struct {
	client *containerd.Client
}

// NewContainerdJobTracker creates a new ContainerdJobTracker instance with the given containerd address.
func NewContainerdJobTracker(containerdAddr string) (*ContainerdJobTracker, error) {
	client, err := containerd.New(containerdAddr)
	if err != nil {
		return nil, err
	}
	return &ContainerdJobTracker{client: client}, nil
}

func (t *ContainerdJobTracker) ListJobs() ([]string, error) {
	ctx := namespaces.WithNamespace(context.Background(), "default")
	containers, err := t.client.Containers(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(containers))
	for i, container := range containers {
		ids[i] = container.ID()
	}
	return ids, nil
}

func (t *ContainerdJobTracker) ListArrayJobs(arrayjobID string) ([]string, error) {
	return helper.ArrayJobID2GUIDs(arrayjobID)
}

func (t *ContainerdJobTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return helper.AddArrayJobAsSingleJobs(jt, t, begin, end, step)
}

func (t *ContainerdJobTracker) JobState(jobID string) (drmaa2interface.JobState, string, error) {
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := t.client.LoadContainer(ctx, jobID)
	if err != nil {
		return drmaa2interface.Undetermined, "can't load container", err
	}
	task, err := container.Task(ctx, nil)
	if err != nil {
		return drmaa2interface.Undetermined, "can't load task", err
	}
	status, err := task.Status(ctx)
	if err != nil {
		return drmaa2interface.Undetermined, "can't get task status", err
	}
	return containerdStatusToDrmaa2State(status), string(status.Status), nil
}

func (t *ContainerdJobTracker) JobInfo(jobID string) (drmaa2interface.JobInfo, error) {
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := t.client.LoadContainer(ctx, jobID)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	info, err := container.Info(ctx)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	status, err := task.Status(ctx)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	image, err := t.client.GetImage(ctx, info.Image)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	return containerdInfoToDRMAA2JobInfo(info, status, image)
}

func containerdInfoToDRMAA2JobInfo(info containers.Container, status containerd.Status, image containerd.Image) (drmaa2interface.JobInfo, error) {
	ji := drmaa2interface.JobInfo{
		ID: info.ID,
		//JobOwner:       info,
		Slots:          1,
		SubmissionTime: info.CreatedAt,
		State:          containerdStatusToDrmaa2State(status),
		ExitStatus:     int(status.ExitStatus),
		//JobCategory:       image.Name(),
		AllocatedMachines: []string{info.ID},
	}

	switch status.Status {
	case containerd.Created, containerd.Running:
		ji.DispatchTime = info.CreatedAt
	case containerd.Stopped:
		ji.FinishTime = status.ExitTime
	}

	return ji, nil
}

func (t *ContainerdJobTracker) JobControl(jobID, action string) error {
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := t.client.LoadContainer(ctx, jobID)
	if err != nil {
		return err
	}
	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}
	switch action {
	case jobtracker.JobControlSuspend:
		return task.Pause(ctx)
	case jobtracker.JobControlResume:
		return task.Resume(ctx)
	case jobtracker.JobControlHold:
		// Not implemented.
		return fmt.Errorf("Hold action is not supported in this implementation")
	case jobtracker.JobControlRelease:
		// Not implemented.
		return fmt.Errorf("Release action is not supported in this implementation")
	case jobtracker.JobControlTerminate:
		return task.Kill(ctx, syscall.SIGKILL)
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

func (t *ContainerdJobTracker) Wait(jobID string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return helper.WaitForState(t, jobID, timeout, state...)
}

func (t *ContainerdJobTracker) DeleteJob(jobID string) error {
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := t.client.LoadContainer(ctx, jobID)
	if err != nil {
		return err
	}
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}

func (t *ContainerdJobTracker) ListJobCategories() ([]string, error) {
	// Not implemented.
	return nil, fmt.Errorf("ListJobCategories is not supported in this implementation")
}

// containerdStatusToDrmaa2State maps a containerd status to a DRMAA2 job state.
func containerdStatusToDrmaa2State(status containerd.Status) drmaa2interface.JobState {
	switch status.Status {
	case containerd.Created:
		return drmaa2interface.Queued
	case containerd.Running:
		return drmaa2interface.Running
	case containerd.Stopped:
		if status.ExitStatus == 0 {
			return drmaa2interface.Done
		}
		return drmaa2interface.Failed
	case containerd.Paused:
		return drmaa2interface.Suspended
	case containerd.Pausing:
		return drmaa2interface.Running
	default:
		return drmaa2interface.Undetermined
	}
}
