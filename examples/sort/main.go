package main

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

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
		JobName:       "sort",
		RemoteCommand: "/usr/bin/sort",
		InputPath:     "/etc/services",
		OutputPath:    "/dev/stdout",
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	job.WaitTerminated(drmaa2interface.InfiniteTime)
}
