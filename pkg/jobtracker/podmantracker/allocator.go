package podmantracker

import (
	"errors"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
)

// PodmanTrackerParams provide parameters which can be passed
// to the SessionManager in order to influence the behviour
// of podman.
type PodmanTrackerParams struct {
	ConnectionURIOverride string
	DisableImagePull      bool
}

type allocator struct{}

func NewAllocator() *allocator {
	return &allocator{}
}

// New is called by the SessionManager when a new JobSession is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	if jobTrackerInitParams != nil {
		podmanParams, ok := jobTrackerInitParams.(PodmanTrackerParams)
		if !ok {
			return nil, errors.New("jobTrackerInitParams for podman has not PodmanTrackerParams type")
		}
		return New(jobSessionName, podmanParams)
	}
	return New(jobSessionName, PodmanTrackerParams{})
}
