package main

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"time"
)

func main() {
	sm, err := drmaa2os.NewKubernetesSessionManager("testdb.db")
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
		// JobName must be unique or not set ("").
		//JobName:       "testjob",
		RemoteCommand: "/bin/sh",
		JobCategory:   "golang",
		Args:          []string{"-c", `sleep 5`},
	}

	fmt.Printf("running job sleep 5\n")
	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	fmt.Println("job submitted successfully")

	job.WaitStarted(drmaa2interface.InfiniteTime)
	<-time.After(time.Millisecond * 500)
	fmt.Printf("Job State: %s\n", job.GetState())

	err = job.Terminate()
	if err != nil {
		fmt.Printf("Error during terminating job: %s\n", err)
	} else {
		fmt.Printf("succesfully terminated job %s\n", job.GetID())
	}

	fmt.Printf("job state: %s\n", job.GetState().String())
	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("final job state: %s\n", job.GetState().String())
	ji, _ := job.GetJobInfo()
	fmt.Printf("job info: %v\n", ji)

	jt.Args = []string{"-c", "exit 1"}

	fmt.Printf("running job exit 1\n")
	job, err = js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	fmt.Println("job submitted successfully")

	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("final job state: %s\n", job.GetState().String())
	ji, _ = job.GetJobInfo()
	fmt.Printf("job info: %v\n", ji)

}
