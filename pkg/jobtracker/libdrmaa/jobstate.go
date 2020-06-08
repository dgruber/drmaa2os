package libdrmaa

import (
	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

// ConvertDRMAAStateToDRMAA2State takes a DRMAA v1 state and converts it in a
// DRMAA2 job state.
func ConvertDRMAAStateToDRMAA2State(pt drmaa.PsType) drmaa2interface.JobState {
	switch pt {
	case drmaa.PsUndetermined:
		return drmaa2interface.Undetermined
	case drmaa.PsQueuedActive:
		return drmaa2interface.Queued
	case drmaa.PsSystemOnHold:
		return drmaa2interface.QueuedHeld
	case drmaa.PsUserOnHold:
		return drmaa2interface.QueuedHeld
	case drmaa.PsUserSystemOnHold:
		return drmaa2interface.QueuedHeld
	case drmaa.PsRunning:
		return drmaa2interface.Running
	case drmaa.PsSystemSuspended:
		return drmaa2interface.Suspended
	case drmaa.PsUserSuspended:
		return drmaa2interface.Suspended
	case drmaa.PsUserSystemSuspended:
		return drmaa2interface.Suspended
	case drmaa.PsDone:
		return drmaa2interface.Done
	case drmaa.PsFailed:
		return drmaa2interface.Failed
	}
	return drmaa2interface.Undetermined
}
