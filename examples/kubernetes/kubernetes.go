package main

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
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
		Args:          []string{"-c", "sleep", "0"},
	}

	for i := 1; i < 100; i++ {
		fmt.Printf("running job %d\n", i)
		job, err := js.RunJob(jt)
		if err != nil {
			panic(err)
		}
		fmt.Println("job submitted successfully")
		err = job.WaitTerminated(drmaa2interface.InfiniteTime)
		if err != nil {
			panic(err)
		}
	}
}
