package dockertracker

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
)

func runJob(jobsession string, cli *client.Client, jt drmaa2interface.JobTemplate) (string, error) {
	if err := checkJobTemplate(jt); err != nil {
		return "", err
	}
	// stage image
	// https://docs.docker.com/engine/api/#api-example
	// https://github.com/moby/moby/blob/master/api/types/container/config.go
	config, err := jobTemplateToContainerConfig(jobsession, jt)
	if err != nil {
		return "", err
	}

	hostConfig, err := jobTemplateToHostConfig(jt)
	if err != nil {
		return "", fmt.Errorf("Docker Host Config: %s", err.Error())
	}

	networkingConfig, err := jobTemplateToNetworkingConfig(jt)
	if err != nil {
		return "", fmt.Errorf("Docker Network Config: %s", err.Error())
	}

	// pull image -> requires internet access
	//_, err = dt.cli.ImagePull(context.Background(), jt.JobCategory, types.ImagePullOptions{})
	// if err != nil {
	//	return "", fmt.Errorf("Error while pulling image: %s", err.Error())
	// }
	ccBody, err := cli.ContainerCreate(context.Background(),
		config,
		hostConfig,
		networkingConfig,
		&v1.Platform{
			Architecture: "amd64",
			OS:           "linux",
		},
		jt.JobName)

	if err != nil {
		return "", fmt.Errorf("creating container: %s", err.Error())
	}

	err = cli.ContainerStart(context.Background(), ccBody.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", fmt.Errorf("starting container: %s", err.Error())
	}

	if jt.OutputPath != "" || jt.ErrorPath != "" {
		stdout := false
		stderr := false

		if jt.OutputPath != "" {
			stdout = true
		}
		if jt.ErrorPath != "" {
			stderr = true
		}

		handleInputOutput(cli,
			ccBody.ID,
			types.ContainerAttachOptions{Stream: true, Stdout: stdout, Stderr: stderr, Logs: true},
			jt.OutputPath,
			jt.ErrorPath)
	}
	return ccBody.ID, nil
}

func handleInputOutput(cli *client.Client, id string, options types.ContainerAttachOptions, stdoutfile, stderrfile string) {
	res, err := cli.ContainerAttach(context.Background(), id, options)
	if err != nil {
		panic(err)
	}
	if stdoutfile == "" {
		stdoutfile = "/dev/null"
	}
	if stderrfile == "" {
		stderrfile = "/dev/null"
	}
	redirectOut(res, stdoutfile, stderrfile)
}

func redirectOut(res types.HijackedResponse, outfilename, errfilename string) {
	go func() {
		var err error

		var outfile *os.File
		if outfilename == "/dev/stdout" {
			outfile = os.Stdout
		} else {
			outfile, err = os.Create(outfilename)
			if err != nil {
				panic(err)
			}
		}

		var errfile *os.File
		if errfilename == "/dev/stderr" {
			errfile = os.Stderr
		} else {
			errfile, err = os.Create(errfilename)
			if err != nil {
				panic(err)
			}
		}
		stdcopy.StdCopy(outfile, errfile, res.Reader)
		outfile.Close()
		errfile.Close()
		res.Close()
	}()
}
