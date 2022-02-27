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

// JobControl action arguments
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
	// ListJobs returns all visible job IDs or an error.
	ListJobs() ([]string, error)
	// ListArrayJobs returns all job IDs an job array ID (or array job ID)
	// represents or an error.
	ListArrayJobs(arrayjobID string) ([]string, error)
	// AddJob typically submits or starts a new job at the backend. The function
	// returns the unique job ID or an error if job submission (or starting of
	// the job in case there is no queueing system) has failed.
	AddJob(jt drmaa2interface.JobTemplate) (string, error)
	// AddArrayJob makes a mass submission of jobs defined by the same job template.
	// Many HPC workload manager support job arrays for submitting 10s of thousands
	// of similar jobs by one call. The additional parameters define how many jobs
	// are submitted by defining a TASK_ID range. Begin is the first task ID (like 1),
	// end is the last task ID (like 10), step is a positive integeger which defines
	// the increments from one task ID to the next task ID (like 1). maxParallel is
	// an arguments representating an optional functionality which instructs the
	// backend to limit maxParallel tasks of this job arary to run in parallel.
	// Note, that jobs use the TASK_ID environment variable to identifiy which
	// task they are and determine that way what to do (like which data set is
	// accessed).
	AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error)
	// JobState returns the DRMAA2 state and substate (free form string) of the job.
	JobState(jobID string) (drmaa2interface.JobState, string, error)
	// JobInfo returns the job status of a job in form of a JobInfo struct or an error.
	JobInfo(jobID string) (drmaa2interface.JobInfo, error)
	// JobControl sends a request to the backend to either "terminate", "suspend",
	// "resume", "hold", or "release" a job. The strings are fixed and are defined
	// by the JobControl constants. This could change in the future to be limited
	// only to constants representing the actions. When the request is not accepted
	// by the system the function must return an error.
	JobControl(jobID, action string) error
	// Wait blocks until the job is either in one of the given states, the max.
	// waiting time (specified by timeout) is reached or an other internal
	// error occured (like job was not found). In case of a timeout also an
	// error must be returned.
	Wait(jobID string, timeout time.Duration, state ...drmaa2interface.JobState) error
	// DeleteJob removes a job from a potential internal database. It does not stop
	// a job. A job must be in an endstate (terminated, failed) in order to call
	// DeleteJob. In case of an error or the job is not in an end state error must be
	// returned. If the backend does not support cleaning up resources for a finished
	// job nil should be returned.
	DeleteJob(jobID string) error
	// ListJobCategories returns a list of job categories which can be used in the
	// JobCategory field of the job template. The list is informational. An example
	// is returning a list of supported container images. AddJob() and AddArrayJob()
	// processes a JobTemplate and hence also the JobCategory field.
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
