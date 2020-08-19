package main

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
)

func main() {
	sm, err := drmaa2os.NewLibDRMAASessionManager("testdb.db")
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
	job, err := js.RunJob(drmaa2interface.JobTemplate{
		JobName:       "job1",
		RemoteCommand: "/bin/sleep",
		Args:          []string{"1"},
	})
	if err != nil {
		fmt.Printf("job submission failed: %v\n", err)
		js.Close()
		os.Exit(1)
	}
	fmt.Printf("job submitted with ID %s\n", job.GetID())
	err = job.WaitStarted(drmaa2interface.InfiniteTime)
	if err != nil {
		fmt.Printf("failed waiting for job to be started: %v\n", err)
		js.Close()
		os.Exit(1)
	}
	fmt.Printf("job started\n")
	jobState := job.GetState()
	fmt.Printf("job state: %s\n", jobState.String())

	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		fmt.Printf("failed waiting for job to be finished: %v\n", err)
		js.Close()
		os.Exit(1)
	}
	fmt.Printf("job finished\n")

	jobState = job.GetState()
	fmt.Printf("job state: %s\n", jobState.String())

	ji, err := job.GetJobInfo()
	ji, err = job.GetJobInfo()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("job info: %v\n", ji)

	js.Close()
	sm.DestroyJobSession("jobsession1")
}
