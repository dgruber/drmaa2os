package containerdtracker

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/helper"
)

func (t *ContainerdJobTracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	ctx := namespaces.WithNamespace(context.Background(), "default")

	if jt.JobName == "" {
		jt.JobName = "drmaa2os-job-" + fmt.Sprintf("%d", time.Now().UnixNano())
	}

	if jt.JobCategory == "" {
		return "",
			fmt.Errorf("JobCategory representing the container image is not set")
	}

	// Pull the image
	image, err := t.client.Pull(ctx, jt.JobCategory, containerd.WithPullUnpack)
	if err != nil {
		return "", fmt.Errorf("Error pulling container image: %v", err)
	}

	labels := make(map[string]string)
	labels["jobSessionName"] = t.JobSessionName
	labels["jobName"] = jt.JobName
	labels["drmaa2"] = "true"
	labels["jobTemplate"], err = helper.JobTemplateToBase64(jt)
	if err != nil {
		return "", fmt.Errorf("Error base64 encoding job template: %v", err)
	}

	jobEnv := make([]string, 0, len(jt.JobEnvironment))
	for k, v := range jt.JobEnvironment {
		jobEnv = append(jobEnv, k+"="+v)
	}

	// Create the container
	container, err := t.client.NewContainer(
		ctx,
		jt.JobName,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(jt.JobName+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image),
			oci.WithProcessArgs(jt.Args...), oci.WithEnv(jobEnv)),
		containerd.WithContainerLabels(labels),
		// io.containerd.runc.v2
		containerd.WithRuntime("remote", nil),
	)
	if err != nil {
		return "", fmt.Errorf("Error creating container: %v", err)
	}

	// Create and start the task (container)
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return "", fmt.Errorf("Error creating container task: %v", err)
	}

	if err := task.Start(ctx); err != nil {
		return "", fmt.Errorf("Error starting container task: %v", err)
	}

	return container.ID(), nil
}
