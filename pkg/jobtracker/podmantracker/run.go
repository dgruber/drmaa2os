package podmantracker

import (
	"context"
	"fmt"

	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/dgruber/drmaa2interface"
)

func RunPodmanContainer(ctx context.Context, jt drmaa2interface.JobTemplate, disablePull bool) (string, error) {
	// context must provide the podman connection: ctx.Value(clientKey).(*Connection)

	spec := specgen.NewSpecGenerator(jt.JobCategory, false)
	spec.Terminal = true

	spec.Command = append([]string{jt.RemoteCommand}, jt.Args...)

	spec.Env = jt.JobEnvironment

	//spec.Mounts = CreateMounts(jt.StageInFiles)

	r, err := containers.CreateWithSpec(ctx, spec, &containers.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return r.ID, containers.Start(ctx, r.ID, &containers.StartOptions{})
}
