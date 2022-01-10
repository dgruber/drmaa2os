package client

import (
	"errors"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
)

type ClientTrackerParams struct {
	// Server of the remote API like "http://localhost:8087"
	Server string
	// Path sets path at server of remote jobtracker API (like "/container")
	Path string
	// Opts are additional settings for the client, like for authentication
	Opts []genclient.ClientOption
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
