package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
)

// Simple example of using drmaa2 with Singularity. Please remove any
// cached instance of the singularity container image so that suspend
// and resume can be demonstrated while pulling the image.

func main() {
	sm, err := drmaa2os.NewSingularitySessionManager(filepath.Join(os.TempDir(), "jobs.db"))
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		js, err = sm.OpenJobSession("jobsession")
		if err != nil {
			panic(err)
		}
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "/bin/sleep",
		Args:          []string{"600"},
		JobCategory:   "shub://GodloveD/lolcow",
		OutputPath:    "/dev/stdout",
		ErrorPath:     "/dev/stderr",
	}
	// set Singularity specific arguments and options
	jt.ExtensionList = map[string]string{
		"debug": "true",
		"pid":   "true",
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	err = job.WaitStarted(drmaa2interface.InfiniteTime)
	if err != nil {
		fmt.Printf("Error while waiting for Singularity container to start (%s).\n", err.Error())
	} else {
		fmt.Printf("Job %s is running.\n", job.GetID())
	}
	<-time.After(time.Second * 10)

	err = job.Suspend()
	if err != nil {
		fmt.Printf("Error while suspending the job (%s).", err.Error())
	} else {
		fmt.Println("Job suspended. Waiting 5 seconds.")
	}
	<-time.After(time.Second * 5)

	err = job.Resume()
	if err != nil {
		fmt.Printf("Error while resuming the job (%s).", err.Error())
	} else {
		fmt.Println("Job resumed.")
	}

	err = job.WaitTerminated(time.Second * 240)
	fmt.Printf("Job is still running: %s\n", err.Error())
	if err := job.Terminate(); err != nil {
		fmt.Printf("Error while terminating Singularity container: %s.\n", err.Error())
	} else {
		fmt.Println("Terminated Singularity container successfully.")
	}

	js.Close()
	sm.DestroyJobSession("jobsession")
}
