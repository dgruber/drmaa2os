package simpletrackerfakes

import (
	"errors"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"time"
)

type FakeJob struct {
	id       string
	session  string
	template drmaa2interface.JobTemplate
	tracker  jobtracker.JobTracker // reference to external job tracker
	jobstate drmaa2interface.JobState
	// config
	faketimeout time.Duration
	endstate    drmaa2interface.JobState
	errorwait   error
}

func NewFakeJob(endstate drmaa2interface.JobState, errMessage string, timeout time.Duration) (fake FakeJob) {
	fake.faketimeout = timeout
	fake.endstate = endstate
	fake.jobstate = endstate
	if errMessage != "" {
		fake.errorwait = errors.New(errMessage)
	}
	return fake
}

func (j FakeJob) GetID() string {
	return j.id
}

func (j FakeJob) GetSessionName() string {
	return j.session
}

func (j FakeJob) GetJobTemplate() (drmaa2interface.JobTemplate, error) {
	return j.template, nil
}

func (j FakeJob) GetState() drmaa2interface.JobState {
	return j.jobstate
}

func (j FakeJob) GetJobInfo() (drmaa2interface.JobInfo, error) {
	return drmaa2interface.JobInfo{}, nil
}

func (j FakeJob) Suspend() error {
	j.jobstate = drmaa2interface.Suspended
	return nil
}

func (j FakeJob) Resume() error {
	j.jobstate = drmaa2interface.Running
	return nil
}

func (j FakeJob) Hold() error {
	j.jobstate = drmaa2interface.QueuedHeld
	return nil
}

func (j FakeJob) Release() error {
	j.jobstate = drmaa2interface.Queued
	return nil
}

func (j FakeJob) Terminate() error {
	j.jobstate = drmaa2interface.Failed
	return nil
}

func (j FakeJob) WaitStarted(timeout time.Duration) error {
	if j.errorwait != nil {
		return j.errorwait
	}
	time.Sleep(j.faketimeout)
	return nil
}

func (j FakeJob) WaitTerminated(timeout time.Duration) error {
	if j.errorwait != nil {
		return j.errorwait
	}
	time.Sleep(j.faketimeout)
	return nil
}

func (j FakeJob) Reap() error {
	return nil
}
