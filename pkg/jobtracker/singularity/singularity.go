package singularity

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

// init registers the singularity tracker at the SessionManager
func init() {
	drmaa2os.RegisterJobTracker(drmaa2os.SingularitySession, NewAllocator())
}

func NewAllocator() *allocator {
	return &allocator{}
}

type allocator struct{}

// New is called by the SessionManager when a new JobSession is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	return New(jobSessionName)
}

// Tracker tracks singularity container.
type Tracker struct {
	processTracker  *simpletracker.JobTracker
	singularityPath string
}

// New creates a new Tracker for Singularity containers.
func New(jobsession string) (*Tracker, error) {
	singularityPath, err := exec.LookPath("singularity")
	if err != nil {
		return nil, fmt.Errorf("singularity command is not found")
	}
	return &Tracker{
		processTracker:  simpletracker.New(jobsession),
		singularityPath: singularityPath,
	}, nil
}

// ListJobs shows all Singularity containers running.
func (dt *Tracker) ListJobs() ([]string, error) {
	return dt.processTracker.ListJobs()
}

// AddJob creates a new Singularity container.
func (dt *Tracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	if jt.JobCategory == "" {
		return "", fmt.Errorf("Singularity container image not specified")
	}
	return dt.processTracker.AddJob(createProcessJobTemplate(jt))
}

// AddArrayJob creates ~(end - begin)/step Singularity containers.
func (dt *Tracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin, end, step, maxParallel int) (string, error) {
	return dt.processTracker.AddArrayJob(createProcessJobTemplate(jt), begin, end, step, maxParallel)
}

// ListArrayJobs shows all containers which belong to a certain job array.
func (dt *Tracker) ListArrayJobs(ID string) ([]string, error) {
	return dt.processTracker.ListArrayJobs(ID)
}

// JobState returns the state of the Singularity container.
func (dt *Tracker) JobState(jobid string) (drmaa2interface.JobState, string, error) {
	return dt.processTracker.JobState(jobid)
}

// JobInfo returns detailed information about the job.
func (dt *Tracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	return dt.processTracker.JobInfo(jobid)
}

// JobControl suspends, resumes, or stops a Singularity container.
func (dt *Tracker) JobControl(jobid, state string) error {
	return dt.processTracker.JobControl(jobid, state)
}

// Wait blocks until either one of the given states is reached or when the timeout occurs.
func (dt *Tracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	return dt.processTracker.Wait(jobid, timeout, state...)
}

// ListJobCategories returns nothing.
func (dt *Tracker) ListJobCategories() ([]string, error) {
	return []string{}, nil
}

// DeleteJob removes the job from the internal storage. It errors
// when the job is not yet in any end state.
func (dt *Tracker) DeleteJob(jobid string) error {
	return dt.processTracker.DeleteJob(jobid)
}
