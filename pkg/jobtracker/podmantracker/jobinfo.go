package podmantracker

import (
	"context"
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/dgruber/drmaa2interface"
)

func ContainerInfo(ctx context.Context, id string) (drmaa2interface.JobInfo, error) {
	c, err := containers.Inspect(ctx, id, nil)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}
	var ji drmaa2interface.JobInfo

	ji.State, ji.SubState, _ = GetContainerState(ctx, id)
	ji.ExitStatus = int(c.State.ExitCode)
	ji.DispatchTime = c.Created
	ji.SubmissionTime = c.State.StartedAt
	ji.FinishTime = c.State.FinishedAt
	ji.Annotation = c.ProcessLabel
	if ji.SubmissionTime != ji.FinishTime {
		ji.WallclockTime = ji.FinishTime.Sub(ji.SubmissionTime)
	}
	ji.JobOwner = c.Config.User
	//ji.SubmissionMachine
	//ji.AllocatedMachines
	ji.ID = c.ID
	ji.QueueName = ""
	ji.Slots = int64(c.HostConfig.CpuCount)
	ji.TerminatingSignal = fmt.Sprintf("%d", c.Config.StopSignal)

	return ji, nil
}
