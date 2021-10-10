package libdrmaa

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/helper"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

// init registers the libdrmaa tracker at the SessionManager
func init() {
	drmaa2os.RegisterJobTracker(drmaa2os.LibDRMAASession, NewAllocator())
}

type allocator struct{}

func NewAllocator() *allocator {
	return &allocator{}
}

// New is called by the SessionManager when a new JobSession is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	return NewDRMAATrackerWithParams(jobTrackerInitParams)
}

// LibDRMAASessionParams contains arguments which can be evaluated
// during DRMAA2 job session creation.
type LibDRMAASessionParams struct {
	// ContactString is required also for opening job sessions
	// hence do not change the name "ContactString" as SessionManager
	// depends on that through reflection
	ContactString string
	// UsePersistentJobStorage saves job ids in a DB file
	// so that they are availabe after an application restart
	// (could be slower with massive amounts of jobs)
	UsePersistentJobStorage bool
	// DBFilePath points to an existing or non-existing boltdb file
	// which is used when persistent storage is used.
	DBFilePath string
}

func (l *LibDRMAASessionParams) SetContact(contact string) {
	l.ContactString = contact
}

// WorkloadManagerType is related to a specific drmaa.so backend as
// there are minor differences in terms of capabilities
type WorkloadManagerType int

const (
	// UnivaGridEngine as recogized drmaa.so backend
	UnivaGridEngine WorkloadManagerType = iota
	// SonOfGridEngine as recogized drmaa.so backend
	SonOfGridEngine
)

// DRMAATracker implements the JobTracker interface with drmaa.so
// as backend for job management. That allows to user drmaa.so
// with a DRMAA2 compatible interface.
type DRMAATracker struct {
	sync.Mutex
	workloadManager WorkloadManagerType
	session         *drmaa.Session
	store           simpletracker.JobStorer
}

func NewDRMAATrackerWithParams(params interface{}) (*DRMAATracker, error) {
	if params == nil {
		return NewDRMAATracker()
	}

	drmaaParams, ok := params.(LibDRMAASessionParams)
	if ok == false {
		return nil, fmt.Errorf("can not initialized DRMAA job tracker as params is not of type LibDRMAASessionParams")
	}

	if drmaaParams.ContactString != "" {
		fmt.Printf("using contact string within drmaa tracker: %s\n", drmaaParams.ContactString)
	}
	if drmaaParams.UsePersistentJobStorage && drmaaParams.DBFilePath == "" {
		return nil,
			fmt.Errorf("using persistent job storage for drmaa session but DBFilePath is not set")
	}

	var session drmaa.Session
	if err := session.Init(drmaaParams.ContactString); err != nil {
		return nil, err
	}
	tracker, err := createTracker(session)
	if err != nil {
		return tracker, err
	}
	if drmaaParams.UsePersistentJobStorage {
		tracker.store, err = simpletracker.NewPersistentJobStore(drmaaParams.DBFilePath)
		if err != nil {
			return nil, err
		}
	}
	return tracker, nil
}

// NewDRMAATracker creates a new JobTracker interface implementation
// which manages jobs through the drmaa (version 1) interface.
func NewDRMAATracker() (*DRMAATracker, error) {
	s, err := drmaa.MakeSession()
	if err != nil {
		return nil, err
	}
	return createTracker(s)
}

func createTracker(s drmaa.Session) (*DRMAATracker, error) {
	// (contact string something like "session=d1b18d34bb44.3871.1722668764")
	// differentiate between different workload manager supporing drmaa.so
	drm, err := s.GetDrmSystem()
	if err != nil {
		return nil, err
	}
	var wlm WorkloadManagerType
	if strings.HasPrefix(drm, "SGE 8.1.") {
		// Son of Grid Engine returns "SGE 8.1.9"
		wlm = SonOfGridEngine
	} else {
		wlm = UnivaGridEngine
	}
	return &DRMAATracker{
		session:         &s,
		store:           simpletracker.NewJobStore(),
		workloadManager: wlm,
	}, nil
}

// DestroySession is not part of the interface but neccessary for
// shutting down the connection to the workload manager.
func (t *DRMAATracker) DestroySession() error {
	return t.session.Exit()
}

// ListJobs returns all jobs previously submitted and still locally cached.
func (t *DRMAATracker) ListJobs() ([]string, error) {
	t.Lock()
	defer t.Unlock()
	return t.store.GetJobIDs(), nil
}

// AddJob makes a new job submission through the underlying drmaa.so
// RunJob function.
func (t *DRMAATracker) AddJob(template drmaa2interface.JobTemplate) (string, error) {
	t.Lock()
	defer t.Unlock()
	jt, err := t.session.AllocateJobTemplate()
	if err != nil {
		return "", err
	}
	defer t.session.DeleteJobTemplate(&jt)

	// a job name might be required even not set in the JobTemplate
	jt.SetJobName("cdrmaatrackerjob")

	err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(template, &jt)
	if err != nil {
		return "", err
	}
	jobID, err := t.session.RunJob(&jt)
	t.store.SaveJob(jobID, template, 0)
	return jobID, err
}

// AddArrayJob submits an array job through the underlying drmaa.so
// RunBulkJobs function.
func (t *DRMAATracker) AddArrayJob(template drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	t.Lock()
	defer t.Unlock()
	jt, err := t.session.AllocateJobTemplate()
	if err != nil {
		return "", err
	}
	defer t.session.DeleteJobTemplate(&jt)
	err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(template, &jt)
	if err != nil {
		return "", err
	}
	taskIDs, err := t.session.RunBulkJobs(&jt, begin, end, step)
	pids := make([]int, 0, 16)
	for range taskIDs {
		pids = append(pids, 0)
	}
	arrayJobID := helper.Guids2ArrayJobID(taskIDs)
	t.store.SaveArrayJob(arrayJobID, pids, template, begin, end, step)
	return arrayJobID, err
}

// ListArrayJobs returns all job IDs of the job array task.
func (t *DRMAATracker) ListArrayJobs(arrayJobID string) ([]string, error) {
	return helper.ArrayJobID2GUIDs(arrayJobID)
}

// JobState returns the current state of the given job.
func (t *DRMAATracker) JobState(jobID string) (drmaa2interface.JobState, string, error) {
	if t == nil || t.session == nil {
		return drmaa2interface.Undetermined, "", fmt.Errorf("no active job session")
	}
	ps, err := t.session.JobPs(jobID)
	if err != nil {
		return drmaa2interface.Undetermined, "", err
	}
	return ConvertDRMAAStateToDRMAA2State(ps), "", nil
}

// JobInfo returns more detailed information about a job when the job is finished.
func (t *DRMAATracker) JobInfo(jobID string) (drmaa2interface.JobInfo, error) {
	// we get the job info when the job is finished - we can also
	// use the DRM system specific calls (like on GE)
	state, _, err := t.JobState(jobID)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}
	if state == drmaa2interface.Failed || state == drmaa2interface.Done {
		// job is in end state
		jinfo, err := t.session.Wait(jobID, 60)
		if err != nil {
			return drmaa2interface.JobInfo{}, err
		}
		return ConvertDRMAAJobInfoToDRMAA2JobInfo(&jinfo), nil
	}
	return drmaa2interface.JobInfo{}, nil
}

// JobControl allows the job to be executed.
func (t *DRMAATracker) JobControl(jobID, action string) error {
	if t == nil || t.session == nil {
		return fmt.Errorf("no active job session")
	}
	switch action {
	case "suspend":
		return t.session.SuspendJob(jobID)
	case "resume":
		return t.session.ResumeJob(jobID)
	case "hold":
		return t.session.HoldJob(jobID)
	case "release":
		return t.session.ReleaseJob(jobID)
	case "terminate":
		return t.session.TerminateJob(jobID)
	}
	return fmt.Errorf("internal: unknown job state change request: %s", action)
}

// Wait blocks until the job reached one of the given states or the timeout is reached.
func (t *DRMAATracker) Wait(jobid string, timeout time.Duration, state ...drmaa2interface.JobState) error {
	// TODO optimize here in case we need wait only for job end states
	return helper.WaitForState(t, jobid, timeout, state...)
}

// DeleteJob removes a job from the internal DB. It can only be removed
// when it is in an end state (failed or done.
func (t *DRMAATracker) DeleteJob(jobID string) error {
	t.Lock()
	defer t.Unlock()
	// job needs to be in an end state
	ps, err := t.session.JobPs(jobID)
	if err != nil {
		return err
	}
	if ps != drmaa.PsDone && ps != drmaa.PsFailed {
		return fmt.Errorf("job is not in an end state (%v)", ps)
	}
	t.store.RemoveJob(jobID)
	return nil
}

// ListJobCategories returns the job categories available at the workload manager.
// Since this is not a drmaa v1 concept we ignore it for now.
func (t *DRMAATracker) ListJobCategories() ([]string, error) {
	return []string{}, nil
}

// Contact() returns the contact string. Implements ContactStringer interface.
// Used for getting the job session name of DRMAA1 out of Grid Engine.
func (t *DRMAATracker) Contact() (string, error) {
	return drmaa.GetContact()
}

// JobTemplate returns the stored job template of the job. This job tracker
// implements the JobTemplater interface additional to the JobTracker interface.
func (jt *DRMAATracker) JobTemplate(jobID string) (drmaa2interface.JobTemplate, error) {
	return jt.store.GetJobTemplate(jobID)
}
