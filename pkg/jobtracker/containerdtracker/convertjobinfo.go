package containerdtracker

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/dgruber/drmaa2interface"
)

const DRMAA2_JobAnnotation = "drmaa2.jobannotation"

func ConvertContainerToInfo(ctx context.Context, c containerd.Container) (drmaa2interface.JobInfo, error) {
	ji := drmaa2interface.JobInfo{
		ID: c.ID(),
	}
	spec, err := c.Spec(ctx)
	if err != nil {
		return ji, fmt.Errorf("could not get spec of container %s: %s", c.ID(), err)
	}
	ji.JobOwner = spec.Process.User.Username
	ji.AllocatedMachines = []string{spec.Hostname}
	ji.Annotation = spec.Annotations[DRMAA2_JobAnnotation]
	info, err := c.Info(ctx)
	ji.SubmissionTime = info.CreatedAt
	ji.DispatchTime = info.CreatedAt

	return ji, nil
}
