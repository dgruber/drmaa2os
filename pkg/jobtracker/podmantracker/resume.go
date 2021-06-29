package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func ResumeContainer(ctx context.Context, id string) error {
	return containers.Unpause(ctx, id, nil)
}
