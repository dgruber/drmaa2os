package main

import (
	"fmt"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

func main() {

	sm, err := drmaa2os.NewDefaultSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.OpenJobSession("jobsession1")
	if err != nil {
		js, err = sm.CreateJobSession("jobsession1", "")
		if err != nil {
			panic(err)
		}
	}
	defer js.Close()

	jt := drmaa2interface.JobTemplate{
		JobName:       "job1",
		RemoteCommand: "sleep",
		Args:          []string{"360"},
	}

	fmt.Printf("Submitting 10000 jobs in an array job and run max. 2 in parallel\n")
	arrayJob, err := js.RunBulkJobs(jt, 1, 10000, 1, 2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Terminate array job\n")
	err = arrayJob.Terminate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wait for all jobs to be finished\n")
	jobs, _ := js.GetJobs(drmaa2interface.CreateJobInfo())

	fmt.Printf("Number of tasks: %d\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("Waiting for task %s to be finished (process ID: %s)\n",
			job.GetID(), job.GetState())
		err = job.WaitTerminated(drmaa2interface.InfiniteTime)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("All jobs deleted\n")

	Benchmark(drmaa2interface.JobTemplate{
		RemoteCommand: "/bin/sh",
		Args:          []string{"-c", `sleep 0.1`},
	}, js, 100)
}

func Benchmark(jt drmaa2interface.JobTemplate, js drmaa2interface.JobSession, tasks int) {

	for i := 1; i <= 20; i += 1 {
		start := time.Now()
		fmt.Printf("Submitting %d jobs in an array job and run max. %d in parallel\n", tasks, i)
		arrayJob, err := js.RunBulkJobs(jt, 1, tasks, 1, i)
		if err != nil {
			panic(err)
		}
		for _, job := range arrayJob.GetJobs() {
			err = job.WaitTerminated(drmaa2interface.InfiniteTime)
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("Took: %s\n", time.Since(start))
	}

}
