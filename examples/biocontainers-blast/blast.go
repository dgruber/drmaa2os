package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
)

//go:embed blast.sh
var script string

func main() {

	// the output should appear in the output subdirectory
	// of the current directory
	cwd, err := GetCwdAndCreateOutputDirectory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// blast job template
	blast := drmaa2interface.JobTemplate{
		JobCategory:      "biocontainers/blast:2.2.31",
		RemoteCommand:    "/bin/bash",
		Args:             []string{"-c", script},
		OutputPath:       "/dev/stdout",
		ErrorPath:        "/dev/stderr",
		WorkingDirectory: "/tmp",
		// mount shared volume
		StageInFiles: map[string]string{
			// local dir: container dir
			cwd + "/output": "/host",
		},
	}

	// create session DB
	file, err := ioutil.TempFile("", "drmaa2os")
	if err != nil {
		panic(err)
	}
	file.Close()

	sm, err := drmaa2os.NewDockerSessionManager(file.Name())
	if err != nil {
		panic(err)
	}
	jobSession, err := sm.OpenJobSession("blast")
	if err != nil {
		jobSession, err = sm.CreateJobSession("blast", "")
		if err != nil {
			panic(err)
		}
	}
	RunJob(jobSession, blast)
	
	sm.DestroyJobSession("blast")
}

func RunJob(js drmaa2interface.JobSession, jt drmaa2interface.JobTemplate) {
	job, err := js.RunJob(jt)
	if err != nil {
		fmt.Printf("failed to submit job: %v\n", err)
		os.Exit(1)
	}
	// wait for job to finish
	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		fmt.Printf("error waiting for job to finish: %v\n", err)
		os.Exit(1)
	}
}

func GetCwdAndCreateOutputDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outputDir := cwd + "/output"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", err
	}
	return cwd, nil
}
