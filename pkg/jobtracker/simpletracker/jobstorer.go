package simpletracker

import (
	"github.com/dgruber/drmaa2interface"
)

// JobStorer has all methods required for storing job related information.
type JobStorer interface {
	SaveJob(jobid string, t drmaa2interface.JobTemplate, pid int)
	HasJob(jobid string) bool
	RemoveJob(jobid string)
	IsArrayJob(jobid string) bool
	SaveArrayJob(arrayjobid string, pids []int, t drmaa2interface.JobTemplate, begin, end, step int)
	SaveArrayJobPID(arrayjobid string, taskid, pid int) error
	GetPID(jobid string) (int, error)
	GetJobIDs() []string
	GetArrayJobTaskIDs(arrayjobID string) []string
	// NewJobID returns a new unique job ID
	NewJobID() string
	// Require job template
	GetJobTemplate(jobid string) (drmaa2interface.JobTemplate, error)
}
