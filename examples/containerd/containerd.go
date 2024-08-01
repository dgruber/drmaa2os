package main

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/containerdtracker"
)

func main() {
	params := containerdtracker.ContainerdTrackerParams{
		// using lima on macOS you need the containerd socket forwared to the host
		// see https://github.com/lima-vm/lima/discussions/1275
		ContainerdAddr: "/Users/daniel/.lima/default/sock/docker.sock",
	}
	sm, err := drmaa2os.NewContainerdSessionManager(
		params, "testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession1", "")
	if err != nil {
		js, err = sm.OpenJobSession("jobsession1")
		if err != nil {
			panic(err)
		}
	}

	jt := drmaa2interface.JobTemplate{
		//JobName:       "testjob",
		JobCategory:   "docker.io/library/busybox:latest",
		RemoteCommand: "/bin/sh",
		Args:          []string{"-c", "echo hello"},
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	job.WaitTerminated(drmaa2interface.InfiniteTime)
}
