package server

import (
	"github.com/dgruber/drmaa2interface"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
)

func ConvertJobInfo(in drmaa2interface.JobInfo) genserver.JobInfo {
	return genserver.JobInfo{
		Id:                in.ID,
		ExitStatus:        in.ExitStatus,
		TerminatingSignal: in.TerminatingSignal,
		Annotation:        in.Annotation,
		State:             string(ConvertJobState(in.State.String())),
		SubState:          in.SubState,
		AllocatedMachines: in.AllocatedMachines,
		SubmissionMachine: in.SubmissionMachine,
		JobOwner:          in.JobOwner,
		Slots:             int(in.Slots),
		QueueName:         in.QueueName,
		WallclockTime:     int64(in.WallclockTime.Seconds()),
		CpuTime:           in.CPUTime,
		SubmissionTime:    in.SubmissionTime,
		DispatchTime:      in.DispatchTime,
		FinishTime:        in.FinishTime,
	}
}

func ConvertJobTemplateToDRMAA2(in genserver.JobTemplate) drmaa2interface.JobTemplate {
	return drmaa2interface.JobTemplate{
		AccountingID:      in.AccountingID,
		Args:              in.Args,
		CandidateMachines: in.CandidateMachines,
		DeadlineTime:      in.DeadlineTime,
		Email:             in.Email,
		EmailOnStarted:    in.EmailOnStarted,
		EmailOnTerminated: in.EmailOnTerminated,
		ErrorPath:         in.ErrorPath,
		InputPath:         in.InputPath,
		JobCategory:       in.JobCategory,
		JobEnvironment:    in.JobEnvironment.AdditionalProperties,
		JobName:           in.JobName,
		JoinFiles:         in.JoinFiles,
		MachineArch:       in.MachineArch,
		MachineOs:         in.MachineOs,
		MaxSlots:          in.MaxSlots,
		MinPhysMemory:     in.MinPhysMemory,
		MinSlots:          in.MinSlots,
		OutputPath:        in.OutputPath,
		Priority:          in.Priority,
		QueueName:         in.QueueName,
		RemoteCommand:     in.RemoteCommand,
		ReRunnable:        in.ReRunnable,
		ReservationID:     in.ReservationID,
		ResourceLimits:    in.ResourceLimits.AdditionalProperties,
		StageInFiles:      in.StageInFiles.AdditionalProperties,
		StageOutFiles:     in.StageOutFiles.AdditionalProperties,
		StartTime:         in.StartTime,
		SubmitAsHold:      in.SubmitAsHold,
		WorkingDirectory:  in.WorkingDirectory,
	}
}

func ConvertJobState(in string) genserver.JobState {
	switch in {
	case drmaa2interface.Done.String():
		return genserver.JobStateDone
	case drmaa2interface.Failed.String():
		return genserver.JobStateFailed
	case drmaa2interface.Queued.String():
		return genserver.JobStateQueued
	case drmaa2interface.QueuedHeld.String():
		return genserver.JobStateQueuedHeld
	case drmaa2interface.Requeued.String():
		return genserver.JobStateRequeued
	case drmaa2interface.RequeuedHeld.String():
		return genserver.JobStateRequeuedHeld
	case drmaa2interface.Running.String():
		return genserver.JobStateRunning
	case drmaa2interface.Suspended.String():
		return genserver.JobStateSuspended
	case drmaa2interface.Undetermined.String():
		return genserver.JobStateUndetermined
	case drmaa2interface.Unset.String():
		return genserver.JobStateUnset
	}
	return genserver.JobStateUndetermined
}
