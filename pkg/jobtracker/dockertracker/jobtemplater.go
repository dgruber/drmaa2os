package dockertracker

import (
	"encoding/json"
	"fmt"

	"encoding/base64"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// JobTemplate returns the JobTemplate for the given jobID. This implements
// the JobTemplater interface for the DockerTracker.
func (dt *DockerTracker) JobTemplate(jobID string) (drmaa2interface.JobTemplate, error) {
	return ReadJobTemplateFromLabel(jobID)
}

// ReadJobTemplateFromLabel reads the "drmaa2jobtemplate" label from
// the specified container. Then it decodes the base64/json encoded
// JobTemplate and returns it. Fo encoding see jobTemplateToContainerConfig().
func ReadJobTemplateFromLabel(containerID string) (drmaa2interface.JobTemplate, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return drmaa2interface.JobTemplate{}, err
	}

	// Fetch the container's current configuration
	inspect, err := cli.ContainerInspect(ctx, containerID)
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
