package containerdtracker

import (
	"errors"

	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
)

// init registers the containerd tracker at the SessionManager
func init() {
	var a allocator
	drmaa2os.RegisterJobTracker(drmaa2os.ContainerdSession, &a)
}

type allocator struct{}

// New is called by the SessionManager when a new JobSession is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	containerdParams, ok := jobTrackerInitParams.(ContainerdTrackerParams)
	if !ok {
		return nil, errors.New("jobTrackerInitParams is not of type []string")
	}
	return NewContainerdJobTracker(jobSessionName,
		containerdParams.ContainerdAddr)
}
