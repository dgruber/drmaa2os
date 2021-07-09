package client

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
)

func ToStringArray(in []*string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s != nil {
			out = append(out, *s)
		}
	}
	return out
}

func ConvertJobTemplateToDRMAA2(in genclient.JobTemplate) drmaa2interface.JobTemplate {
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

func ConvertJobTemplate(in drmaa2interface.JobTemplate) genclient.JobTemplate {
	return genclient.JobTemplate{
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
		JobEnvironment:    genclient.JobTemplate_JobEnvironment{AdditionalProperties: in.JobEnvironment},
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
		ResourceLimits:    genclient.JobTemplate_ResourceLimits{AdditionalProperties: in.ResourceLimits},
		StageInFiles:      genclient.JobTemplate_StageInFiles{AdditionalProperties: in.StageInFiles},
		StageOutFiles:     genclient.JobTemplate_StageOutFiles{AdditionalProperties: in.StageOutFiles},
		StartTime:         in.StartTime,
		SubmitAsHold:      in.SubmitAsHold,
		WorkingDirectory:  in.WorkingDirectory,
	}
}

func ConvertJobInfoToDRMAA2(in genclient.JobInfo) drmaa2interface.JobInfo {
	fmt.Printf("JOB INFO: %v\n", in)
	fmt.Printf("state: %v\n", in.State)
	return drmaa2interface.JobInfo{
		ID:                in.Id,
		ExitStatus:        in.ExitStatus,
		TerminatingSignal: in.TerminatingSignal,
		Annotation:        in.Annotation,
		State:             ConvertJobStateToDRMAA2(in.State),
		SubState:          in.SubState,
		AllocatedMachines: in.AllocatedMachines,
		SubmissionMachine: in.SubmissionMachine,
		JobOwner:          in.JobOwner,
		Slots:             int64(in.Slots),
		QueueName:         in.QueueName,
		//WallclockTime:     time.ParseDuration(in.WallclockTime * time.Second),
		CPUTime:        in.CpuTime,
		SubmissionTime: in.SubmissionTime,
		DispatchTime:   in.DispatchTime,
		FinishTime:     in.FinishTime,
	}
}

func ConvertJobStateToDRMAA2(in string) drmaa2interface.JobState {
	switch in {
	case string(genclient.JobStateDone):
		return drmaa2interface.Done
	case string(genclient.JobStateFailed):
		return drmaa2interface.Failed
	case string(genclient.JobStateQueued):
		return drmaa2interface.Queued
	case string(genclient.JobStateQueuedHeld):
		return drmaa2interface.QueuedHeld
	case string(genclient.JobStateRequeued):
		return drmaa2interface.Requeued
	case string(genclient.JobStateRequeuedHeld):
		return drmaa2interface.RequeuedHeld
	case string(genclient.JobStateRunning):
		return drmaa2interface.Running
	case string(genclient.JobStateSuspended):
		return drmaa2interface.Suspended
	case string(genclient.JobStateUndetermined):
		return drmaa2interface.Undetermined
	case string(genclient.JobStateUnset):
		return drmaa2interface.Unset
	}
	return drmaa2interface.Undetermined
}
