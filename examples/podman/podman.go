package main

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	// also needs to register Podman backend when podmantracker package is not accessed
	"github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
)

func main() {
	sm, err := drmaa2os.NewPodmanSessionManager(
		podmantracker.PodmanTrackerParams{
			ConnectionURIOverride: "ssh://vagrant@localhost:2222/tmp/podman.sock?secure=False",
		}, "testdb.db")
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
		RemoteCommand: "/bin/sh",
		JobCategory:   "busybox:latest",
		Args:          []string{"-c", `sleep 2 && ps -ef && ls -lisha && whoami && exit 13`},
		OutputPath:    "/dev/stdout",
		ErrorPath:     "/dev/null",
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	jobs, err := js.GetJobs(drmaa2interface.CreateJobInfo())
	if err != nil {
		fmt.Printf("failed getting jobs: %v\n", err)
	}
	fmt.Printf("found %d jobs\n", len(jobs))

	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		panic(err)
	}

	if job.GetState() == drmaa2interface.Done {
		fmt.Printf("Job %s finished successfully\n", job.GetID())
	} else {
		ji, err := job.GetJobInfo()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Job %s finished with exit code %d\n", job.GetID(), ji.ExitStatus)
	}

	ji, err := job.GetJobInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("JobInfo: %v\n", ji)

	name, _ := js.GetSessionName()
	fmt.Printf("Job session: %s\n", name)

	job.Reap()
	js.Close()
	sm.DestroyJobSession("jobsession1")
}
