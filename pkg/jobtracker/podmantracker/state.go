package podmantracker

import (
	"context"
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/dgruber/drmaa2interface"
)

func GetContainerState(ctx context.Context, id string) (drmaa2interface.JobState, string, error) {
	c, err := containers.Inspect(ctx, id, nil)
	if err != nil {
		return drmaa2interface.Undetermined, "", fmt.Errorf("container with ID %s not found: %v", id, err)
	}
	if c.State.Restarting {
		return drmaa2interface.Running, "restarting", nil
	}
	if c.State.Paused {
		return drmaa2interface.Suspended, "", nil
	}
	if c.State.Running {
		return drmaa2interface.Running, "", nil
	}
	if c.State.ExitCode != 0 {
		return drmaa2interface.Failed, "", nil
	}
	return drmaa2interface.Done, "", nil
}
