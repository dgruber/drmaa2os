package simpletracker

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dgruber/drmaa2interface"
)

// JobEvent is send whenever a job status change is happening
// to inform all registered listeners.
type JobEvent struct {
	JobID    string
	JobState drmaa2interface.JobState
	JobInfo  drmaa2interface.JobInfo
	callback chan bool // if set sends true if event was distributed
}

// PubSub distributes job status change events to clients which
// Register() at PubSub.
type PubSub struct {
	sync.Mutex

	// go routines write into that channel when process has finished
	jobch chan JobEvent

	// maps a jobid to functions registered for waiting for a specific
	// state of that job
	waitFunctions map[string][]waitRequest

	// feed by bookKeeper: current state
	jobState map[string]drmaa2interface.JobState
	jobInfo  map[string]drmaa2interface.JobInfo

	jobstore JobStorer
}

// NewPubSub returns an initialized PubSub structure and
// the JobEvent channel which is used by the caller to publish
// job events (i.e. job state transitions).
func NewPubSub(jobstore JobStorer) (*PubSub, chan JobEvent) {

	jeCh := make(chan JobEvent, 1)

	pubSub := &PubSub{
		jobch:         jeCh,
		waitFunctions: make(map[string][]waitRequest),
		jobState:      make(map[string]drmaa2interface.JobState),
		jobInfo:       make(map[string]drmaa2interface.JobInfo),
	}

	if jobstore != nil {
		pubSub.jobstore = jobstore
		// get all information
		for _, existingJobID := range jobstore.GetJobIDs() {
			jobinfo, err := jobstore.GetJobInfo(existingJobID)
			if err != nil {
				pubSub.jobState[existingJobID] = drmaa2interface.Undetermined
				pubSub.jobInfo[existingJobID] = drmaa2interface.JobInfo{
					ID:    existingJobID,
					State: drmaa2interface.Undetermined,
				}
			} else {
				// restore state from disk - might be wrong
				// running processes might not be running anymore

				pubSub.jobInfo[existingJobID] = jobinfo
				pubSub.jobState[existingJobID] = jobinfo.State

				// check non final states
				if jobinfo.State == drmaa2interface.Running {
					pid, err := jobstore.GetPID(existingJobID)
					if err != nil {
						continue
					}
					running, err := IsPidRunning(pid)
					if err != nil || running == false {
						pubSub.jobState[existingJobID] = drmaa2interface.Undetermined
						jobinfo.State = drmaa2interface.Undetermined
						jobinfo.SubState = "finished before application started"
						pubSub.jobInfo[existingJobID] = jobinfo
					}
				}

				if jobinfo.State == drmaa2interface.Queued {
					// for queued jobs we don't even have a pid
					pubSub.jobState[existingJobID] = drmaa2interface.Undetermined
					jobinfo.State = drmaa2interface.Undetermined
					jobinfo.SubState = "queued before application started"
					pubSub.jobInfo[existingJobID] = jobinfo
				}

			}
		}
	}

	return pubSub, jeCh
}

// Register returns a channel which emits a job state once the given
// job transitions in one of the given states. If job is already
// in the expected state it returns nil as channel and nil as error.
//
// TODO add method for removing specific wait functions.
func (ps *PubSub) Register(jobid string, states ...drmaa2interface.JobState) (chan drmaa2interface.JobState, error) {
	ps.Lock()
	defer ps.Unlock()

	// check if job is already in the expected state
	state, exists := ps.jobState[jobid]
	if exists {
		for _, expectedState := range states {
			if expectedState == state {
				return nil, nil
			}
		}
		if state == drmaa2interface.Failed || state == drmaa2interface.Done {
			return nil, errors.New("job already finished")
		}
	} else {
		return nil, fmt.Errorf("job %s does not exist", jobid)
	}

	waitChannel := make(chan drmaa2interface.JobState, 1)
	ps.waitFunctions[jobid] = append(ps.waitFunctions[jobid],
		waitRequest{ExpectedState: states, WaitChannel: waitChannel})
	return waitChannel, nil
}

// Unregister removes all functions waiting for a specific job and
// all occurences of the job itself.
func (ps *PubSub) Unregister(jobid string) {
	ps.Lock()
	defer ps.Unlock()
	delete(ps.waitFunctions, jobid)
	delete(ps.jobState, jobid)
	delete(ps.jobInfo, jobid)
}

// NotifyAndWait sends a job event and waits until it was distributed
// to all waiting functions.
func (ps *PubSub) NotifyAndWait(evt JobEvent) {
	evt.callback = make(chan bool, 1)
	ps.jobch <- evt
	<-evt.callback
}

// waitRequest defines when to notify (ExpectedState) and where to notify (WaitChannel)
type waitRequest struct {
	ExpectedState []drmaa2interface.JobState
	WaitChannel   chan drmaa2interface.JobState
}

// StartBookKeeper processes all job state changes from the process trackers
// and notifies registered wait functions.
func (ps *PubSub) StartBookKeeper() {
	go func() {
		for event := range ps.jobch {
			ps.Lock()
			ps.jobState[event.JobID] = event.JobState
			if info, exists := ps.jobInfo[event.JobID]; exists {
				ps.jobInfo[event.JobID] = mergeJobInfo(info, event.JobInfo)
			} else {
				// deep copy
				ps.jobInfo[event.JobID] = mergeJobInfo(drmaa2interface.CreateJobInfo(),
					event.JobInfo)
			}
			if ps.jobstore != nil {
				ps.jobstore.SaveJobInfo(event.JobID, ps.jobInfo[event.JobID])
			}
			// inform registered functions
			for _, waiter := range ps.waitFunctions[event.JobID] {
				// inform when expected state is reached
				for i := range waiter.ExpectedState {
					if event.JobState == waiter.ExpectedState[i] {
						waiter.WaitChannel <- event.JobState
					}
				}
			}
			ps.Unlock()
			if event.callback != nil {
				event.callback <- true
			}
		}
	}()
}
func (ps *PubSub) GetJobInfo(jobID string) (drmaa2interface.JobInfo, error) {
	ps.Lock()
	defer ps.Unlock()
	jobInfo, exists := ps.jobInfo[jobID]
	if !exists {
		return drmaa2interface.JobInfo{}, fmt.Errorf("does not exist")
	}
	return mergeJobInfo(drmaa2interface.CreateJobInfo(),
		jobInfo), nil
}

func mergeJobInfo(oldJI, newJI drmaa2interface.JobInfo) drmaa2interface.JobInfo {
	if newJI.ID != "" {
		oldJI.ID = newJI.ID
	}
	if newJI.ExitStatus != drmaa2interface.UnsetNum {
		oldJI.ExitStatus = newJI.ExitStatus
	}
	if newJI.TerminatingSignal != "" {
		oldJI.TerminatingSignal = newJI.TerminatingSignal
	}
	if newJI.Annotation != "" {
		oldJI.Annotation = newJI.Annotation
	}
	if newJI.State != drmaa2interface.Unset {
		oldJI.State = newJI.State
	}
	if newJI.SubState != "" {
		oldJI.SubState = newJI.SubState
	}
	if newJI.AllocatedMachines != nil {
		oldJI.AllocatedMachines = make([]string, 0, len(newJI.AllocatedMachines))
		copy(oldJI.AllocatedMachines, newJI.AllocatedMachines)
	}
	if newJI.SubmissionMachine != "" {
		oldJI.SubmissionMachine = newJI.SubmissionMachine
	}
	if newJI.JobOwner != "" {
		oldJI.JobOwner = newJI.JobOwner
	}
	if newJI.Slots != drmaa2interface.UnsetNum {
		oldJI.Slots = newJI.Slots
	}
	if newJI.QueueName != "" {
		oldJI.QueueName = newJI.QueueName
	}
	if newJI.WallclockTime.Microseconds() > oldJI.WallclockTime.Microseconds() {
		oldJI.WallclockTime = newJI.WallclockTime
	}
	if newJI.CPUTime > oldJI.CPUTime {
		oldJI.CPUTime = newJI.CPUTime
	}
	if !newJI.SubmissionTime.IsZero() {
		oldJI.SubmissionTime = newJI.SubmissionTime
	}
	if !newJI.DispatchTime.IsZero() {
		oldJI.DispatchTime = newJI.DispatchTime
	}
	if !newJI.FinishTime.IsZero() {
		oldJI.FinishTime = newJI.FinishTime
	}
	if newJI.ExtensionList != nil {
		if oldJI.ExtensionList == nil {
			oldJI.ExtensionList = make(map[string]string, len(newJI.ExtensionList))
		}
		for k, v := range newJI.ExtensionList {
			oldJI.ExtensionList[k] = v
		}
	}
	return oldJI
}
