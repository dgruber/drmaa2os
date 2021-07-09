package client

import (
	"errors"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
)

type ClientTrackerParams struct {
	Server string
}

type allocator struct{}

func NewAllocator() *allocator {
	return &allocator{}
}

// New is called by the SessionManager when a new JobSession is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	if jobTrackerInitParams != nil {
		clientTrackerParams, ok := jobTrackerInitParams.(ClientTrackerParams)
		if !ok {
			return nil, errors.New("jobTrackerInitParams for remote client is not of type ClientTrackerParams")
		}
		return New(jobSessionName, clientTrackerParams)
	}
	return New(jobSessionName, ClientTrackerParams{})
}
