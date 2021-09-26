package drmaa2os

import (
	"errors"
	"fmt"
	"log"

	"code.cloudfoundry.org/lager"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"

	"github.com/dgruber/drmaa2os/pkg/storage"
)

// SessionType represents the selected resource manager.
type SessionType int

const (
	// DefaultSession handles jobs as processes
	DefaultSession SessionType = iota
	// DockerSession manages Docker containers
	DockerSession
	// CloudFoundrySession manages Cloud Foundry application tasks
	CloudFoundrySession
	// KubernetesSession creates Kubernetes jobs
	KubernetesSession
	// SingularitySession manages Singularity containers
	SingularitySession
	// SlurmSession manages slurm jobs
	SlurmSession
	// LibDRMAASession manages jobs through libdrmaa.so
	LibDRMAASession
	// PodmanSession manages jobs as podman containers either locally or remote
	PodmanSession
	// RemoteSession manages jobs over the network through a remote server
	RemoteSession
	// ExternalSession can be used by external JobTracker implementations
	// during development time before they get added here
	ExternalSession
)

func init() {
	// initialize job tracker registration map
	atomicTrackers.Store(make(map[SessionType]jobtracker.Allocator))
}

// RegisterJobTracker registers a JobTracker implementation at session manager
// so that it can be used. This is done in the init() method of the JobTracker
// implementation. That means the application which wants to use a specific JobTracker
// needs to import the JobTracker implementation package with _.
//
// Like when Docker needs to be used as job management backend:
//
// import _ "github.com/dgruber/drmaa2os/pkg/jobtracker/pkg/dockertracker"
//
// When multiple backends to be used, all of them needs to be imported so
// that they are registered in the main application.
func RegisterJobTracker(sessionType SessionType, tracker jobtracker.Allocator) {
	trackerMutex.Lock()
	jtMap := atomicTrackers.Load().(map[SessionType]jobtracker.Allocator)
	if jtMap == nil {
		jtMap = make(map[SessionType]jobtracker.Allocator)
	}
	jtMap[sessionType] = tracker
	atomicTrackers.Store(jtMap)
	trackerMutex.Unlock()
}

// SessionManager allows to create, list, and destroy job, reserveration,
// and monitoring sessions. It also returns holds basic information about
// the resource manager and its capabilities.
type SessionManager struct {
	store                  storage.Storer
	log                    lager.Logger
	sessionType            SessionType
	jobTrackerCreateParams interface{}
}

// NewDefaultSessionManager creates a SessionManager which starts jobs
// as processes.
func NewDefaultSessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, DefaultSession)
}

// NewSingularitySessionManager creates a new session manager creating and
// maintaining jobs as Singularity containers.
func NewSingularitySessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, SingularitySession)
}

// NewDockerSessionManager creates a SessionManager which maintains jobs as
// Docker containers.
func NewDockerSessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, DockerSession)
}

// NewCloudFoundrySessionManager creates a SessionManager which maintains jobs
// as Cloud Foundry tasks.
// addr needs to point to the cloud controller API and username and password
// needs to be set as well.
func NewCloudFoundrySessionManager(addr, username, password, dbpath string) (*SessionManager, error) {
	sm, err := makeSessionManager(dbpath, CloudFoundrySession)
	if err != nil {
		return sm, err
	}
	// specific parameters for Cloud Foundry
	sm.jobTrackerCreateParams = []string{addr, username, password}
	return sm, nil
}

// NewKubernetesSessionManager creates a new session manager which uses
// Kubernetes tasks as execution backend for jobs. The first parameter must
// be either a *kubernetes.Clientset or nil to allocate a new one.
func NewKubernetesSessionManager(cs interface{}, dbpath string) (*SessionManager, error) {
	sm, err := makeSessionManager(dbpath, KubernetesSession)
	if err != nil {
		return sm, err
	}
	// when a job session is created is requires a kubernetes clientset
	sm.jobTrackerCreateParams = cs
	return sm, nil
}

// NewSlurmSessionManager creates a new session manager which wraps the
// slurm command line for managing jobs.
func NewSlurmSessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, SlurmSession)
}

// NewLibDRMAASessionManager creates a new session manager which wraps
// libdrmaa.so (DRMAA v1) through the Go DRMAA library. Please check out
// the details of github.com/dgruber/drmaa before using it. Make sure
// all neccessary paths are set (C header files, LD_LIBRARY_PATH).
func NewLibDRMAASessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, LibDRMAASession)
}

// NewLibDRMAASessionManagerWithParams creates a Go DRMAA session manager
// like NewLibDRMAASessionManager but with additional parameters. The
// parameters must be of type _libdrmaa.LibDRMAASessionParams_.
func NewLibDRMAASessionManagerWithParams(ds interface{}, dbpath string) (*SessionManager, error) {
	sm, err := makeSessionManager(dbpath, LibDRMAASession)
	if err != nil {
		return sm, err
	}
	sm.jobTrackerCreateParams = ds
	return sm, nil
}

// NewPodmanSessionManager creates a new session manager for Podman.
// The first parameter is either nil for using defaults or must be
// of type _podmantracker.PodmanTrackerParams_.
func NewPodmanSessionManager(ps interface{}, dbpath string) (*SessionManager, error) {
	sm, err := makeSessionManager(dbpath, PodmanSession)
	if err != nil {
		return sm, err
	}
	// specific parameters for Podman
	sm.jobTrackerCreateParams = ps
	return sm, nil
}

// NewRemoteSessionManager create a new session manager for accessing
// a remote jobtracker server implementation which can be of any
// backend type.
func NewRemoteSessionManager(rs interface{}, dbpath string) (*SessionManager, error) {
	sm, err := makeSessionManager(dbpath, RemoteSession)
	if err != nil {
		return sm, err
	}
	// specific parameters for remote (like server address)
	sm.jobTrackerCreateParams = rs
	return sm, nil
}

// NexExternalSessionManager creates a new external session. This can be
// used when a JobTrack is implemented outside of the repository.
// Note that only one ExternalSession is available so it makes sense to
// add a constant here.
func NexExternalSessionManager(dbpath string) (*SessionManager, error) {
	return makeSessionManager(dbpath, ExternalSession)
}

// CreateJobSession creates a new JobSession for managing jobs.
func (sm *SessionManager) CreateJobSession(name, contact string) (drmaa2interface.JobSession, error) {
	if err := sm.create(storage.JobSessionType, name, contact); err != nil {
		return nil, err
	}
	// allocate a registered job tracker - registration happens
	// when the package is imported in the init method of the
	// JobTracker implementation package
	jt, err := sm.newRegisteredJobTracker(name, sm.jobTrackerCreateParams)
	if err != nil {
		return nil, err
	}
	js := newJobSession(name, []jobtracker.JobTracker{jt})

	// for libdrmaa return contact string and store it for open job session
	if sm.sessionType == LibDRMAASession {
		if contactStringer, ok := jt.(jobtracker.ContactStringer); ok {
			contact, err := contactStringer.Contact()
			if err != nil {
				return nil, fmt.Errorf("Failed to get contact string after session creation: %v", err)
			}
			// store new contact string for job session
			sm.store.Put(storage.JobSessionType, name, contact)
		}
	}

	return js, nil
}

// CreateReservationSession creates a new ReservationSession.
func (sm *SessionManager) CreateReservationSession(name, contact string) (drmaa2interface.ReservationSession, error) {
	return nil, ErrorUnsupportedOperation
}

// OpenMonitoringSession opens a session for monitoring jobs.
func (sm *SessionManager) OpenMonitoringSession(sessionName string) (drmaa2interface.MonitoringSession, error) {
	return nil, errors.New("(TODO) not implemented")
}

// OpenJobSession creates a new session for managing jobs. The semantic of a job session
// and the job session name depends on the resource manager.
func (sm *SessionManager) OpenJobSession(name string) (drmaa2interface.JobSession, error) {
	if exists := sm.store.Exists(storage.JobSessionType, name); !exists {
		return nil, errors.New("JobSession does not exist")
	}

	// require a copy as it gets modified
	createParams := sm.jobTrackerCreateParams

	// restore contact string from storage and set it as ContactString
	// in job tracker create params
	if sm.sessionType == LibDRMAASession {
		contact, err := sm.store.Get(storage.JobSessionType, name)
		if err != nil {
			return nil, fmt.Errorf("could not get contact string for job session: %s: %v",
				name, err)
		}
		log.Printf("using internal DRMAA job session %s with contact string %s\n", name, contact)
		err = TryToSetContactString(createParams, contact)
		if err != nil {
			return nil, fmt.Errorf("could not set new contact string for opening job session %s: %v",
				name, err)
		}
	}

	jt, err := sm.newRegisteredJobTracker(name, createParams)
	if err != nil {
		return nil, err
	}
	js := JobSession{
		name:    name,
		tracker: []jobtracker.JobTracker{jt},
	}
	return &js, nil
}

// OpenReservationSession opens a reservation session.
func (sm *SessionManager) OpenReservationSession(name string) (drmaa2interface.ReservationSession, error) {
	return nil, ErrorUnsupportedOperation
}

// DestroyJobSession destroys a job session by name.
func (sm *SessionManager) DestroyJobSession(name string) error {
	return sm.delete(storage.JobSessionType, name)
}

// DestroyReservationSession removes a reservation session.
func (sm *SessionManager) DestroyReservationSession(name string) error {
	return ErrorUnsupportedOperation
}

// GetJobSessionNames returns a list of all job sessions.
func (sm *SessionManager) GetJobSessionNames() ([]string, error) {
	return sm.store.List(storage.JobSessionType)
}

// GetReservationSessionNames returns a list of all reservation sessions.
func (sm *SessionManager) GetReservationSessionNames() ([]string, error) {
	return nil, ErrorUnsupportedOperation
}

// GetDrmsName returns the name of the distributed resource manager.
func (sm *SessionManager) GetDrmsName() (string, error) {
	return "drmaa2os", nil
}

// GetDrmsVersion returns the version of the distributed resource manager.
func (sm *SessionManager) GetDrmsVersion() (drmaa2interface.Version, error) {
	return drmaa2interface.Version{Minor: "0", Major: "1"}, nil
}

// Supports returns true of false of the given Capability is supported by DRMAA2OS.
func (sm *SessionManager) Supports(capability drmaa2interface.Capability) bool {
	return false
}

// RegisterEventNotification creates an event channel which emits events when
// the conditions described in the given notification specification are met.
func (sm *SessionManager) RegisterEventNotification() (drmaa2interface.EventChannel, error) {
	return nil, ErrorUnsupportedOperation
}
