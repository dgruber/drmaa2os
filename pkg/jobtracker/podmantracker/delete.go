package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func DeleteContainer(ctx context.Context, id string) error {
	f := false
	t := true
	return containers.Remove(ctx, id, &containers.RemoveOptions{
		Force:   &f,
		Volumes: &t,
	})
}
