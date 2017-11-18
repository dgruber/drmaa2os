package fakes

import (
	"errors"
	"github.com/dgruber/drmaa2interface"
	"time"
)

type Job struct {
	ID               string
	Session          string
	Template         drmaa2interface.JobTemplate
	Jobinfo          drmaa2interface.JobInfo
	State            drmaa2interface.JobState
	ErrorWhenSuspend bool
}

func (j *Job) GetID() string {
	return j.ID
}

func (j *Job) GetSessionName() string {
	return j.Session
}

func (j *Job) GetJobTemplate() (drmaa2interface.JobTemplate, error) {
	return j.Template, nil
}

func (j *Job) GetJobInfo() (drmaa2interface.JobInfo, error) {
	return j.Jobinfo, nil
}

func (j *Job) GetState() drmaa2interface.JobState {
	return j.State
}

func (j *Job) Suspend() error {
	if j.ErrorWhenSuspend {
		j.State = drmaa2interface.Running
		return errors.New("Some error happend")
	}
	j.State = drmaa2interface.Suspended
	return nil
}

func (j *Job) Resume() error {
	j.State = drmaa2interface.Running
	return nil
}

func (j *Job) Hold() error {
	j.State = drmaa2interface.QueuedHeld
	return nil
}

func (j *Job) Release() error {
	j.State = drmaa2interface.Running
	return nil
}

func (j *Job) Terminate() error {
	j.State = drmaa2interface.Failed
	return nil
}

func (j *Job) WaitStarted(timeout time.Duration) error {
	j.State = drmaa2interface.Running
	return nil
}

func (j *Job) WaitTerminated(timeout time.Duration) error {
	j.State = drmaa2interface.Failed
	return nil
}

func (j *Job) Reap() error {
	return nil
}
