package jobtracker

import (
	"time"

	"github.com/dgruber/drmaa2interface"
)

// Allocator contains all what is required to create a new JobTacker
// instance. A JobTracker implementation needs to register the Allocator
// implementation in its init method where it needs to call RegisterJobTracker()
// of the drmaa2os SessionManager. The jobTrackerInitParams are an optional
// way for parameterizing the JobTracker creation method.
type Allocator interface {
	New(jobSessionName string, jobTrackerInitParams interface{}) (JobTracker, error)
}

// JobControl arguments
const JobControlTerminate = "terminate"
const JobControlSuspend = "suspend"
const JobControlResume = "resume"
const JobControlHold = "hold"
const JobControlRelease = "release"

// JobTracker is the interface a basic JobTracker needs to implement
// in order to be able to be hooked into the DRMAA2OS framework.
// Additionaly functionalities of a JobTracker are defined by additional
// interfaces implemented by the same object. Those interfaces are
// listed below (ContactStringer, JobTemplater, Closer, Monitorer).
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

// JobTemplater is a JobTracker which can retrieve the JobTemplate of a job.
type JobTemplater interface {
	JobTemplate(jobid string) (drmaa2interface.JobTemplate, error)
}

// Closer is a JobTracker which needs to disengage from the backend when the
// session is closed so that a new JobTracker with using the same session name
// can be created again.
type Closer interface {
	Close() error
}

// Monitorer is a JobTracker which implements the functions required for
// serving the required capabilities for implementing a MonitoringSession.
// Sources of the machines, jobs, and job states can be implemented
// differently as there is no local persistency layer required.
type Monitorer interface {
	OpenMonitoringSession(name string) error
	GetAllJobIDs(filter *drmaa2interface.JobInfo) ([]string, error)
	GetAllQueueNames(filter []string) ([]string, error)
	GetAllMachines(filter []string) ([]drmaa2interface.Machine, error)
	CloseMonitoringSession(name string) error
	// JobInfoFromMonitor might collect job state and job info in a
	// different way as a JobSession with persistent storage does
	JobInfoFromMonitor(id string) (drmaa2interface.JobInfo, error)
}

// constants for Monitorer struct extensions

const DRMAA2_MS_JOBINFO_WORKINGDIR = "workingdir"
const DRMAA2_MS_JOBINFO_COMMANDLINE = "commandline"
const DRMAA2_MS_JOBINFO_JOBCATEGORY = "category"
