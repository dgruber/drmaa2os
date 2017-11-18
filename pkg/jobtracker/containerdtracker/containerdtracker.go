package containerdtracker

import (
	"github.com/dgruber/drmaa2interface"
	"time"
)

type ContainerDTracker struct{}

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

func (dt *ContainerDTracker) JobState(jobid string) drmaa2interface.JobState {
	return drmaa2interface.Undetermined
}

func (dt *ContainerDTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	return drmaa2interface.JobInfo{}, nil
}

func (dt *ContainerDTracker) JobControl(jobid, state string) error {
	return nil
}

func (dt *ContainerDTracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return nil
}

func (dt *ContainerDTracker) DeleteJob(jobid string) error {
	return nil
}
