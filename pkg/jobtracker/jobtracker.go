package jobtracker

import (
	"time"

	"github.com/dgruber/drmaa2interface"
)

// Allocator contains all what is required to create a new JobTacker
// instance. A JobTracker implementation needs to register the Allocator
// implementation in its init method where it needs to call RegisterJobTracker()
// of the drmaa2os SessionManager. The jobTrackerInitParams are an optional
// way for parameterize the JobTracker creation method.
type Allocator interface {
	New(jobSessionName string, jobTrackerInitParams interface{}) (JobTracker, error)
}

// JobControl arguments
const JobControlTerminate = "terminate"
const JobControlSuspend = "suspend"
const JobControlResume = "resume"
const JobControlHold = "hold"
const JobControlRelease = "release"

type JobTracker interface {
	ListJobs() ([]string, error)
	ListArrayJobs(string) ([]string, error)
	AddJob(jt drmaa2interface.JobTemplate) (string, error)
	AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error)
	JobState(jobid string) (drmaa2interface.JobState, string, error)
	JobInfo(jobid string) (drmaa2interface.JobInfo, error)
	JobControl(jobid, state string) error
	Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error
	DeleteJob(jobid string) error
	ListJobCategories() ([]string, error)
}

// ContactStringer is a JobTracker which offers the Contact() method which
// returns the contact string. Used in the DRMAA1 JobTracker.
type ContactStringer interface {
	Contact() (string, error)
}
