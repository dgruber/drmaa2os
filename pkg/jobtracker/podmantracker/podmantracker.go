package podmantracker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/helper"
)

// PodmanTracker implements the JobTracker interface for managing
// containers as jobs.
type PodmanTracker struct {
	// connectionContext must have a connection to Podman stored
	connectionContext context.Context
	// disableImagePull if set to true prevents that images gets pulled
	disableImagePull bool
}

// init registers the Podman tracker at the SessionManager
func init() {
	// TODO add PodmanSession
	drmaa2os.RegisterJobTracker(drmaa2os.PodmanSession, NewAllocator())
}

// New creates a new connection to Podman and returns a JobTracker interface
// for Podman. If connectionURIOverride is set it uses this URI for the
// new Podman connection otherwise the connection is established through
// the socket.
//
// According to Podman:
//   "A valid URI connection should be scheme://
//    For example tcp://localhost:<port>
//    or unix:///run/podman/podman.sock
//    or ssh://<user>@<host>[:port]/run/podman/podman.sock?secure=True"
func New(jobSessionName string, params PodmanTrackerParams) (*PodmanTracker, error) {
	if params.ConnectionURIOverride == "" {
		sock_dir := os.Getenv("XDG_RUNTIME_DIR")
		params.ConnectionURIOverride = "unix:" + sock_dir + "/podman/podman.sock"
	}

	connText, err := bindings.NewConnection(context.Background(), params.ConnectionURIOverride)
	if err != nil {
		return nil, fmt.Errorf("could not create connection to podman: %v", err)
	}

	return &PodmanTracker{
		connectionContext: connText,
		disableImagePull:  params.DisableImagePull,
	}, nil
}

func (p *PodmanTracker) ListJobs() ([]string, error) {
	return ListPodmanContainers(p.connectionContext)
}

func (p *PodmanTracker) AddJob(template drmaa2interface.JobTemplate) (string, error) {
	return RunPodmanContainer(p.connectionContext, template, p.disableImagePull)
}

func (p *PodmanTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return helper.AddArrayJobAsSingleJobs(jt, p, begin, end, step)
}

func (p *PodmanTracker) ListArrayJobs(arrayjobid string) ([]string, error) {
	return helper.ArrayJobID2GUIDs(arrayjobid)
}

func (p *PodmanTracker) JobState(jobid string) (drmaa2interface.JobState, string, error) {
	return GetContainerState(p.connectionContext, jobid)
}

func (p *PodmanTracker) JobInfo(jobid string) (drmaa2interface.JobInfo, error) {
	return ContainerInfo(p.connectionContext, jobid)
}

func (p *PodmanTracker) JobControl(jobid, action string) error {
	if p == nil {
		return fmt.Errorf("no active job session")
	}
	switch action {
	case "suspend":
		return PauseContainer(p.connectionContext, jobid)
	case "resume":
		return ResumeContainer(p.connectionContext, jobid)
	case "hold":
		return fmt.Errorf("hold is not implemented as there is no queueing")
	case "release":
		return fmt.Errorf("release is not implemented as there is no queueing")
	case "terminate":
		return TerminateContainer(p.connectionContext, jobid)
	}
	return fmt.Errorf("internal: unknown job state change request: %s", action)
}

// Wait until the job has a certain DRMAA2 state or return an error if the state
// is unreachable.
func (p *PodmanTracker) Wait(jobid string, timeout time.Duration, states ...drmaa2interface.JobState) error {
	// this can be replaced with an event based API if available
	return helper.WaitForState(p, jobid, timeout, states...)
}

// DeleteJob removes the container and its volumes from the node. The container
// must be in an end state (i.e. not running anymore).
func (p *PodmanTracker) DeleteJob(jobid string) error {
	return DeleteContainer(p.connectionContext, jobid)
}

// ListJobCategories returns all localy available container images which can
// be used in JobCategory of the JobTemplate.
func (p *PodmanTracker) ListJobCategories() ([]string, error) {
	return ListContainerImages(p.connectionContext)
}
