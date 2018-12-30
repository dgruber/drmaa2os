package main

import (
	"os"
	"path/filepath"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
)

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
		RemoteCommand: "/bin/sh",
		Args:          []string{"-c", "sleep 1 && echo container task: $TASK_ID"},
		JobCategory:   "shub://GodloveD/lolcow",
		OutputPath:    "/dev/stdout",
		ErrorPath:     "/dev/stderr",
	}
	// set Singularity specific arguments and options
	jt.ExtensionList = map[string]string{
		"pid": "true",
	}

	// mass submit of 1000 singularity containers echo running with
	// a different TASK_ID environment variable - throttling to have
	// max. 10 containers running at the same point in time
	jobarray, err := js.RunBulkJobs(jt, 1, 100, 1, 10)
	if err != nil {
		panic(err)
	}
	jobs := jobarray.GetJobs()
	for _, job := range jobs {
		job.WaitTerminated(drmaa2interface.InfiniteTime)
	}

	js.Close()
	sm.DestroyJobSession("jobsession")
}
