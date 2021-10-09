package main

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	// need to register process tracker either by importing it with _
	// or when parameters are used just importing the package
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

func main() {
	params := simpletracker.SimpleTrackerInitParams{
		// note that enabling persistent storage for
		// job IDs reduces performance at least
		// by a factor of 10 (100ms per job vs 8ms)
		// as each job causes DB interaction.
		PersistentStorage:   false,
		PersistentStorageDB: "job.db",
	}
	sm, err := drmaa2os.NewDefaultSessionManagerWithParams(
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
