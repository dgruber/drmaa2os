package simpletracker

import (
	"github.com/dgruber/drmaa2interface"
)

// InternalJob represents a process as a job.
type InternalJob struct {
	TaskID int
	State  drmaa2interface.JobState
	PID    int
}
