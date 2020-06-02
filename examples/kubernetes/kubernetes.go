package main

import (
	"fmt"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
)

func createJobSession(sm drmaa2interface.SessionManager) drmaa2interface.JobSession {
	js, err := sm.CreateJobSession("jobsession1", "")
	if err != nil {
		js, err = sm.OpenJobSession("jobsession1")
		if err != nil {
			panic(err)
		}
	}
	return js
}

func print(ji drmaa2interface.JobInfo) {
	fmt.Printf("Submission time: %s\n", ji.SubmissionMachine)
	fmt.Printf("Dispatch time: %s\n", ji.DispatchTime)
	fmt.Printf("End time: %s\n", ji.FinishTime)
	fmt.Printf("State: %s\n", ji.State)
	fmt.Printf("Job ID: %s\n", ji.ID)
}

func main() {
	sm, err := drmaa2os.NewKubernetesSessionManager(nil, "testdb.db")
	if err != nil {
		panic(err)
	}

	js := createJobSession(sm)
	defer js.Close()

	jt := drmaa2interface.JobTemplate{
		// JobName must be unique or not set ("").
		// JobName:       "testjob",
		RemoteCommand: "/bin/sh",
		JobCategory:   "golang",
		Args:          []string{"-c", `sleep 2`},
	}

	fmt.Printf("running job \"sleep 1\"\n")

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	fmt.Println("job submitted successfully / waiting until finished")

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	ji, err := job.GetJobInfo()
	if err != nil {
		panic(err)
	}
	print(ji)

	fmt.Println("Starting job array with 1000 jobs")
	fmt.Println(time.Now())
	jobs, err := js.RunBulkJobs(jt, 1, 100, 1, 100)
	if err != nil {
		panic(err)
	}
	for _, j := range jobs.GetJobs() {
		j.WaitTerminated(drmaa2interface.InfiniteTime)
		fmt.Printf("Job %s finished\n", j.GetID())
	}
	fmt.Println("All finished")
	fmt.Println(time.Now())
	for _, j := range jobs.GetJobs() {
		j.Reap()
	}
	fmt.Println(time.Now())
}
