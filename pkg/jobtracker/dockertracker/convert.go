package dockertracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

func checkJobTemplate(jt drmaa2interface.JobTemplate) error {
	if jt.JobCategory == "" {
		return errors.New("JobCategory must be set to container image name")
	}
	return nil
}

func setEnv(env map[string]string) []string {
	if env == nil {
		return nil
	}
	envList := make([]string, 0, len(env))
	for key, value := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}
	return envList
}

func newPortSet(ports string) nat.PortSet {
	if ports == "" {
		return nil
	}
	portSet, _, err := nat.ParsePortSpecs(strings.Split(ports, ","))
	if err != nil {
		return nil
	}
	return portSet
}

func newPortBindings(ports string) nat.PortMap {
	if ports == "" {
		return nil
	}
	_, portMap, err := nat.ParsePortSpecs(strings.Split(ports, ","))
	if err != nil {
		return nil
	}
	return portMap
}

// https://github.com/moby/moby/blob/master/api/types/container/config.go
func jobTemplateToContainerConfig(jobsession string, jt drmaa2interface.JobTemplate) (*container.Config, error) {
	var cc container.Config

	cc.Labels = map[string]string{"drmaa2_jobsession": jobsession}
	cc.WorkingDir = jt.WorkingDirectory
	cc.Image = jt.JobCategory

	if len(jt.CandidateMachines) == 1 {
		cc.Hostname = jt.CandidateMachines[0]
	}

	if jt.RemoteCommand != "" {
		cmdSlice := strslice.StrSlice{jt.RemoteCommand}
		cmdSlice = append(cmdSlice, jt.Args...)
		cc.Cmd = cmdSlice
	}

	if jt.WorkingDirectory != "" {
		cc.WorkingDir = jt.WorkingDirectory
	}

	cc.Env = setEnv(jt.JobEnvironment)
	// Docker specific settings in the extensions
	if jt.ExtensionList != nil {
		cc.User = jt.ExtensionList["user"]
		cc.ExposedPorts = newPortSet(jt.ExtensionList["exposedPorts"])

		if value, found := jt.ExtensionList["entrypoint"]; found && jt.ExtensionList["entrypoint"] != "" {
			cc.Entrypoint = strings.Split(value, " ")
		}
	}

	//cc.Tty = true // merges stderr into stdout
	cc.AttachStdout = true
	cc.AttachStderr = true

	// TODO extensions
	// cc.Volumes

	return &cc, nil
}

func jobTemplateToHostConfig(jt drmaa2interface.JobTemplate) (*container.HostConfig, error) {
	var hc container.HostConfig
	//hc.CpusetMems
	//hc.Ulimits
	for outer, inner := range jt.StageInFiles {
		// TODO: Mapping must be the other way around as you might
		// want to map the same file/directory to multiple locations
		// inside the container. So the destination is unique and
		// must be the key (and the source the value).
		hc.Binds = append(hc.Binds, fmt.Sprintf("%s:%s", outer, inner))
	}
	if jt.ExtensionList != nil {
		restart, exists := jt.ExtensionList["restart"]
		if exists {
			// experimental / should be handled outside
			// like --restart=unless-stopped
			hc.RestartPolicy = container.RestartPolicy{
				Name: restart,
			}
		}
		privileged, exists := jt.ExtensionList["privileged"]
		if exists && strings.ToUpper(privileged) != "FALSE" {
			hc.Privileged = true
		}
		net, exists := jt.ExtensionList["net"]
		if exists {
			// like --net host
			hc.NetworkMode = container.NetworkMode(net)
		}
		ipc, exists := jt.ExtensionList["ipc"]
		if exists {
			hc.IpcMode = container.IpcMode(ipc)
		}
		uts, exists := jt.ExtensionList["uts"]
		if exists {
			hc.UTSMode = container.UTSMode(uts)
		}
		pid, exists := jt.ExtensionList["pid"]
		if exists {
			hc.PidMode = container.PidMode(pid)
		}
		rm, exists := jt.ExtensionList["rm"]
		if exists {
			if strings.ToUpper(rm) == "TRUE" {
				hc.AutoRemove = true
			}
		}
	}

	hc.PortBindings = newPortBindings(jt.ExtensionList["exposedPorts"])
	return &hc, nil
}

func jobTemplateToNetworkingConfig(jt drmaa2interface.JobTemplate) (*network.NetworkingConfig, error) {
	var nw network.NetworkingConfig
	// extensions
	return &nw, nil
}

func containersToJobList(jobsession string, containers []types.Container) []string {
	out := make([]string, 0, len(containers))
	for _, c := range containers {
		if js, exists := c.Labels["drmaa2_jobsession"]; exists && js == jobsession {
			out = append(out, c.ID)
		}
	}
	return out
}

func containerToDRMAA2State(state *types.ContainerState) drmaa2interface.JobState {
	// Status be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
	if state.Status == "exited" {
		if state.ExitCode != 0 {
			return drmaa2interface.Failed
		} else {
			return drmaa2interface.Done
		}
	}
	if state.OOMKilled {
		return drmaa2interface.Failed
	}
	if state.Dead {
		if state.ExitCode != 0 {
			return drmaa2interface.Failed
		} else {
			return drmaa2interface.Done
		}
	}
	if state.Paused {
		return drmaa2interface.Suspended
	}
	if state.Restarting {
		return drmaa2interface.Queued
	}
	if state.Running {
		return drmaa2interface.Running
	}
	return drmaa2interface.Undetermined
}

func containerToDRMAA2JobInfo(c types.ContainerJSON) (ji drmaa2interface.JobInfo, err error) {
	ji.ID = c.ID
	ji.Slots = 1
	if c.Config != nil {
		ji.AllocatedMachines = []string{c.Config.Hostname}
	}
	if c.State != nil {
		ji.ExitStatus = c.State.ExitCode
		finished, err := time.Parse(time.RFC3339Nano, c.State.FinishedAt)
		if err == nil {
			ji.FinishTime = finished
		}
		started, err := time.Parse(time.RFC3339Nano, c.State.StartedAt)
		if err == nil {
			ji.DispatchTime = started
		}
		ji.State = containerToDRMAA2State(c.State)
	}
	submitted, err := time.Parse(time.RFC3339Nano, c.Created)
	if err == nil {
		ji.SubmissionTime = submitted
	}
	ji.JobOwner = c.Config.User

	// put more details about the job in the extensions
	ji.ExtensionList = map[string]string{}
	ji.ExtensionList[jobtracker.DRMAA2_MS_JOBINFO_WORKINGDIR] = c.Config.WorkingDir
	ji.ExtensionList[jobtracker.DRMAA2_MS_JOBINFO_COMMANDLINE] = strings.Join(c.Config.Cmd, " ")
	ji.ExtensionList[jobtracker.DRMAA2_MS_JOBINFO_JOBCATEGORY] = c.Config.Image

	return ji, nil
}

func arrayJobID2GUIDs(id string) ([]string, error) {
	var guids []string
	err := json.Unmarshal([]byte(id), &guids)
	if err != nil {
		return nil, err
	}
	return guids, nil
}

func guids2ArrayJobID(guids []string) string {
	id, err := json.Marshal(guids)
	if err != nil {
		return ""
	}
	return string(id)
}

func isInExpectedState(state drmaa2interface.JobState, states ...drmaa2interface.JobState) bool {
	for _, expectedState := range states {
		if state == expectedState {
			return true
		}
	}
	return false
}
