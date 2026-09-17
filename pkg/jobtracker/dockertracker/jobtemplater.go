package dockertracker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/client"
)

// JobTemplate returns the JobTemplate for the given jobID. This implements
// the JobTemplater interface for the DockerTracker.
func (dt *DockerTracker) JobTemplate(jobID string) (drmaa2interface.JobTemplate, error) {
	if err := dt.check(); err != nil {
		return drmaa2interface.JobTemplate{}, err
	}
	return readJobTemplateFromLabel(dt.cli, jobID)
}

// ReadJobTemplateFromLabel reads the "drmaa2jobtemplate" label from
// the specified container. Then it decodes the base64/json encoded
// JobTemplate and returns it. For encoding see jobTemplateToContainerConfig().
// It creates its own Docker client; DockerTracker.JobTemplate reuses the
// client of the tracker instead.
func ReadJobTemplateFromLabel(containerID string) (drmaa2interface.JobTemplate, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return drmaa2interface.JobTemplate{}, err
	}
	// Close only releases idle connections and always returns nil
	defer cli.Close()
	return readJobTemplateFromLabel(cli, containerID)
}

func readJobTemplateFromLabel(cli *client.Client, containerID string) (drmaa2interface.JobTemplate, error) {
	// Fetch the container's current configuration
	inspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return drmaa2interface.JobTemplate{}, err
	}

	// Read the "template" label
	value, ok := inspect.Config.Labels[ContainerLabelJobTemplate]
	if !ok {
		return drmaa2interface.JobTemplate{}, fmt.Errorf("label 'drmaa2jobtemplate' not found")
	}

	// base64 decode the drmaa2interface.JobTemplate
	decodedTemplate, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return drmaa2interface.JobTemplate{}, err
	}
	var jobTemplate drmaa2interface.JobTemplate
	err = json.Unmarshal(decodedTemplate, &jobTemplate)
	if err != nil {
		return drmaa2interface.JobTemplate{}, err
	}

	return jobTemplate, nil
}
