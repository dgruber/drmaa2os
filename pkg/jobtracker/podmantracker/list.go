package podmantracker

import (
	"context"

	"github.com/containers/podman/v3/pkg/bindings/containers"
)

func ListPodmanContainers(ctx context.Context) ([]string, error) {
	containerList, err := containers.List(ctx, &containers.ListOptions{})
	if err != nil {
		return nil, err
	}
	containers := make([]string, 0, 16)
	for _, cl := range containerList {
		containers = append(containers, cl.ID)
	}
	return containers, nil
}
