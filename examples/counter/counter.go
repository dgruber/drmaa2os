package main

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	// need to register process tracker
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

func main() {
	sm, err := drmaa2os.NewDefaultSessionManager("testdb.db")
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
		JobName:       "testjob",
		RemoteCommand: "./plus.sh",
		InputPath:     "in.txt",
		OutputPath:    "out.txt",
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	job.WaitTerminated(drmaa2interface.InfiniteTime)

	for i := 2; i < 1000; i++ {
		if i%2 == 0 {
			jt.InputPath = "out.txt"
			jt.OutputPath = "in.txt"
		} else {
			jt.InputPath = "in.txt"
			jt.OutputPath = "out.txt"
		}
		job, err = js.RunJob(jt)
		if err != nil {
			panic(err)
		}
		err = job.WaitTerminated(drmaa2interface.InfiniteTime)
		if err != nil {
			panic(err)
		}
	}
}
