package podmantracker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/containers/podman/v3/libpod/network/types"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/dgruber/drmaa2interface"
)

// RunPodmanContainer converts a DRMAA2 job template into a container spec and
// runs the container with podman.
//
// The context must provide the podman connection: ctx.Value(clientKey).(*Connection)
func RunPodmanContainer(ctx context.Context, jt drmaa2interface.JobTemplate, disablePull bool) (string, error) {
	spec, err := CreateContainerSpec(jt)
	if err != nil {
		return "", err
	}
	r, err := containers.CreateWithSpec(ctx, spec, &containers.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if jt.OutputPath == "" && jt.ErrorPath == "" {
		return r.ID, containers.Start(ctx, r.ID, &containers.StartOptions{})
	}

	// if stdout and stderr is set attach to container
	err = containers.Start(ctx, r.ID, &containers.StartOptions{})
	if err != nil {
		return r.ID, err
	}

	var stdoutCh chan string = nil
	var stderrCh chan string = nil

	useStdout := false
	stdout, stdoutOpened, err := setWriterOrNot(jt.OutputPath)
	if err != nil {
		return "", err
	}
	if stdout != nil {
		stdoutCh = make(chan string, 512)
		useStdout = true
	}
	useStderr := false
	stderr, stderrOpened, err := setWriterOrNot(jt.ErrorPath)
	if err != nil {
		return "", err
	}
	if stderr != nil {
		stderrCh = make(chan string, 512)
		useStderr = true
	}
	t := true
	go func() {
		err = containers.Logs(ctx, r.ID, &containers.LogOptions{Follow: &t, Stderr: &useStderr, Stdout: &useStdout},
			stdoutCh, stderrCh)
		if err != nil {
			fmt.Printf("failed attaching to logs: %v\n", err)
		}
	}()
	if useStdout {
		go func() {
			for line := range stdoutCh {
				fmt.Fprintf(stdout, "%s\n", line)
			}
			if stdoutOpened {
				stdout.Close()
			}
		}()
	}

	if useStderr {
		go func() {
			for line := range stderrCh {
				fmt.Fprintf(stderr, "%s\n", line)
			}
			if stderrOpened {
				stderr.Close()
			}
		}()
	}

	return r.ID, nil
}

func setWriterOrNot(path string) (*os.File, bool, error) {
	if path == "" {
		return nil, false, nil
	}
	if path == "/dev/stdout" {
		return os.Stdout, false, nil
	}
	if path == "/dev/stderr" {
		return os.Stderr, false, nil
	}
	if path == os.DevNull {
		return nil, false, nil
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, false, fmt.Errorf("failed to append container output to file %s: %v", path, err)
	}
	return file, true, nil
}

func CreateContainerSpec(jt drmaa2interface.JobTemplate) (*specgen.SpecGenerator, error) {
	spec := specgen.NewSpecGenerator(jt.JobCategory, false)

	spec.Terminal = true
	spec.Command = append([]string{jt.RemoteCommand}, jt.Args...)
	spec.Env = jt.JobEnvironment

	// CandidateMachines could be also used for remote ssh based
	// Podman invocation unlike Docker...
	if len(jt.CandidateMachines) > 0 {
		if len(jt.CandidateMachines) == 1 {
			spec.Hostname = jt.CandidateMachines[0]
		} else {
			return nil, fmt.Errorf("CandidateMachines in JobTemplate should have max. 1 entry but has %d",
				len(jt.CandidateMachines))
		}
	}
	if jt.WorkingDirectory != "" {
		spec.WorkDir = jt.WorkingDirectory
	}

	if value, exists := hasExtension(jt, "user"); exists {
		spec.User = value
	}

	if value, exists := hasExtension(jt, "exposedPorts"); exists {
		mappings := make([]types.PortMapping, 0, len(strings.Split(value, ",")))
		for _, portspair := range strings.Split(value, ",") {
			// portspair is format [hostip:]hostport:containerport
			// when hostip is missing we add a -:
			if strings.Count(portspair, ":") == 1 {
				portspair = "-:" + portspair
			}
			ports := strings.Split(portspair, ":")
			if len(ports) != 3 {
				return nil, fmt.Errorf("ports extension should be of format [hostIP:]hostPort:containerPort,hostPort:containerPort,... but contains %s", portspair)
			}
			containerPort, err := strconv.Atoi(ports[2])
			if err != nil {
				return nil, fmt.Errorf("container port is not a number in ports JobTemplate extension")
			}
			hostPort, err := strconv.Atoi(ports[1])
			if err != nil {
				return nil, fmt.Errorf("host port is not a number in ports JobTemplate extension")
			}
			hostIP := ""
			if ports[0] != "-" {
				hostIP = ports[0]
			}
			mappings = append(mappings, types.PortMapping{
				HostIP:        hostIP,
				ContainerPort: uint16(containerPort),
				HostPort:      uint16(hostPort),
			})
		}
		spec.PortMappings = mappings
	}

	if value, exists := hasExtension(jt, "privileged"); exists && strings.ToUpper(value) == "TRUE" {
		spec.Privileged = true
	}

	if value, exists := hasExtension(jt, "restart"); exists {
		spec.RestartPolicy = value
	}

	if value, exists := hasExtension(jt, "ipc"); exists {
		var err error
		spec.IpcNS, err = specgen.ParseNamespace(value)
		if err != nil {
			return nil, fmt.Errorf("failed to set ipc namespace from job template: %v", err)
		}
	}

	if value, exists := hasExtension(jt, "uts"); exists {
		var err error
		spec.UtsNS, err = specgen.ParseNamespace(value)
		if err != nil {
			return nil, fmt.Errorf("failed to set uts namespace from job template: %v", err)
		}
	}

	if value, exists := hasExtension(jt, "pid"); exists {
		var err error
		spec.PidNS, err = specgen.ParseNamespace(value)
		if err != nil {
			return nil, fmt.Errorf("failed to set pid namespace from job template: %v", err)
		}
	}

	if value, exists := hasExtension(jt, "rm"); exists && strings.ToUpper(value) == "TRUE" {
		spec.Remove = true
	}

	//spec.Mounts = CreateMounts(jt.StageInFiles)
	return spec, nil
}

func hasExtension(jt drmaa2interface.JobTemplate, extension string) (string, bool) {
	if jt.ExtensionList == nil {
		return "", false
	}
	value, exists := jt.ExtensionList[extension]
	return value, exists
}
