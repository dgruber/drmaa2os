package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func TerminateContainer(ctx context.Context, id string) error {
	return containers.Kill(ctx, id, nil)
}
