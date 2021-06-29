package main

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
	// need to register Podman backend when podmantracker package is not accessed
	//_ "github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
)

func main() {
	sm, err := drmaa2os.NewPodmanSessionManager(podmantracker.PodmanTrackerParams{
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
		RemoteCommand: "sleep",
		JobCategory:   "busybox:latest",
		Args:          []string{"10"},
	}

	_, err = js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	jobs, _ := js.GetJobs(drmaa2interface.CreateJobInfo())
	j, err := js.WaitAnyTerminated(jobs, drmaa2interface.InfiniteTime)

	if j.GetState() == drmaa2interface.Done {
		fmt.Printf("Job %s finished successfully\n", j.GetID())
	} else {
		fmt.Printf("Job %s finished with failure\n", j.GetID())
	}

	ji, err := j.GetJobInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("JobInfo: %v\n", ji)

	name, _ := js.GetSessionName()
	fmt.Printf("Job session: %s\n", name)

	j.Reap()
	js.Close()
	sm.DestroyJobSession("jobsession1")
}
