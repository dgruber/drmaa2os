package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func PauseContainer(ctx context.Context, id string) error {
	return containers.Pause(ctx, id, nil)
}
