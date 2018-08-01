package simpletrackerfakes

import (
	"errors"
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"time"
)

type JobTracker struct {
	jobsession string
	jobs       map[string]drmaa2interface.Job
	state      map[string]drmaa2interface.JobState
	info       map[string]drmaa2interface.JobInfo
	lastjobid  int
}

func New(sessionname string) *JobTracker {
	return &JobTracker{
		jobsession: sessionname,
		jobs:       make(map[string]drmaa2interface.Job),
		state:      make(map[string]drmaa2interface.JobState),
		info:       make(map[string]drmaa2interface.JobInfo),
		lastjobid:  0,
	}
}

func (jt *JobTracker) ListJobs() ([]string, error) {
	var jobs []string
	for id, _ := range jt.jobs {
		jobs = append(jobs, id)
	}
	return jobs, nil
}

func (jt *JobTracker) ListJobCategories() ([]string, error) {
	return []string{"image", "otherimage"}, nil
}

func (jt *JobTracker) AddJob(t drmaa2interface.JobTemplate) (string, error) {
	jt.lastjobid++
	jobid := fmt.Sprintf("%d", jt.lastjobid)

	job := FakeJob{
		id:          jobid,
		session:     jt.jobsession,
		template:    t,
		tracker:     jt,
		faketimeout: time.Millisecond * 100,
	}

	jt.jobs[jobid] = job
	jt.state[jobid] = drmaa2interface.Running
	jt.info[jobid] = drmaa2interface.JobInfo{
		ID:                jobid,
		State:             drmaa2interface.Running,
		AllocatedMachines: []string{"localhost"},
		SubmissionMachine: "localhost",
		JobOwner:          "testuser",
		Slots:             1,
		SubmissionTime:    time.Now(),
		DispatchTime:      time.Now(),
	}

	return jobid, nil
}

func (jt *JobTracker) AddArrayJob(t drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	jt.lastjobid++
	jobid := fmt.Sprintf("%d", jt.lastjobid)

	for i := begin; i <= end; i += step {
		arrayjobid := fmt.Sprintf("%s.%d", jobid, i)
		job := FakeJob{
			id:       arrayjobid,
			session:  jt.jobsession,
			template: t,
			tracker:  jt,
		}

		jt.jobs[arrayjobid] = job
		jt.state[arrayjobid] = drmaa2interface.Running
	}

	return jobid, nil
}

func (jt *JobTracker) ListArrayJobs(string) ([]string, error) {
	return nil, nil
}

func (jt *JobTracker) JobState(jobid string) drmaa2interface.JobState {
	return jt.state[jobid]
}

func (jt *JobTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	jinfo, exists := jt.info[jobid]
	if exists == false {
		return drmaa2interface.CreateJobInfo(), errors.New("job does not exist")
	}
	return jinfo, nil
}

func (jt *JobTracker) JobControl(jobid, state string) error {
	switch state {
	case "suspend":
		jt.state[jobid] = drmaa2interface.Suspended
	case "resume":
		jt.state[jobid] = drmaa2interface.Running
	case "hold":
		jt.state[jobid] = drmaa2interface.QueuedHeld
	case "release":
		jt.state[jobid] = drmaa2interface.Running
	case "terminate":
		jt.state[jobid] = drmaa2interface.Failed
	}
	return nil
}

func (jt *JobTracker) Wait(jobid string, d time.Duration, states ...drmaa2interface.JobState) error {
	jt.state[jobid] = states[0]
	return nil
}

func (jt *JobTracker) DeleteJob(jobid string) error {
	delete(jt.jobs, jobid)
	delete(jt.state, jobid)
	delete(jt.info, jobid)
	return nil
}
