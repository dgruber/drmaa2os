package main

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"time"
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
		JobName:       "job",
		RemoteCommand: "sleep",
		Args:          []string{"0"},
	}

	var jis []drmaa2interface.JobInfo

	fmt.Println("Start running 1000 sleeper jobs sequentially")
	start := time.Now()

	for i := 0; i < 1000; i++ {
		if job, err := js.RunJob(jt); err != nil {
			panic(err)
		} else {
			err = job.WaitTerminated(drmaa2interface.InfiniteTime)
			if err != nil {
				panic(err)
			}
			if ji, err := job.GetJobInfo(); err != nil {
				panic(err)
			} else {
				jis = append(jis, ji)
			}
		}
	}

	fmt.Printf("It took: %s\n", time.Since(start).String())

	fmt.Printf("Job Info: %v\n", jis[0])

	fmt.Println("Start running 1000 sleeper jobs in parallel")

	start = time.Now()
	var jobs []drmaa2interface.Job

	for i := 0; i < 1000; i++ {
		if job, err := js.RunJob(jt); err != nil {
			panic(err)
		} else {
			jobs = append(jobs, job)
		}
	}

	for i := 0; i < 1000; i++ {
		jobs[i].WaitTerminated(drmaa2interface.InfiniteTime)
	}

	fmt.Printf("It took: %s\n", time.Since(start).String())

}
