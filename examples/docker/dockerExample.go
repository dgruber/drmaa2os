package main

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	// need to register docker backend
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
)

func removeJob(jobs []drmaa2interface.Job, job drmaa2interface.Job) (result []drmaa2interface.Job) {
	if job == nil {
		return jobs
	}
	for i := 0; i < len(jobs); i++ {
		if job.GetID() != jobs[i].GetID() {
			result = append(result, jobs[i])
		}
	}
	return result
}

func main() {
	sm, err := drmaa2os.NewDockerSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession1", "docker")
	if err != nil {
		js, err = sm.OpenJobSession("jobsession1")
		if err != nil {
			panic(err)
		}
	}

	jt := drmaa2interface.JobTemplate{
		// Job names must be unique
		JobName:       "job1",
		RemoteCommand: "sleep",
		JobCategory:   "dgruber/hello",
		Args:          []string{"20"},
	}

	job1, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	jt.JobName = "job2"
	job2, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	jobs, _ := js.GetJobs(drmaa2interface.CreateJobInfo())
	for i := 0; i < 2; i++ {
		j, err := js.WaitAnyTerminated(jobs, drmaa2interface.InfiniteTime)
		jobs = removeJob(jobs, j)

		if err != nil {
			fmt.Printf("Error while waiting for jobs to finish: %s\n", err.Error())
			break
		}
		if j.GetState() == drmaa2interface.Done {
			fmt.Printf("Job %s finished successfully\n", j.GetID())
		} else {
			fmt.Printf("Job %s finished with failure\n", j.GetID())
		}
	}

	jt.JobName = "job3"
	job3, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	jt.JobName = "job4"
	job4, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job3.WaitTerminated(drmaa2interface.InfiniteTime)
	if _, err := job2.GetJobInfo(); err != nil {
		panic(err)
	}

	job4.WaitTerminated(drmaa2interface.InfiniteTime)

	name, _ := js.GetSessionName()
	fmt.Printf("Job session: %s\n", name)

	// we need to delete the container as the container name must be unique
	job1.Reap()
	job2.Reap()
	job3.Reap()
	job4.Reap()

	js.Close()
	sm.DestroyJobSession("jobsession1")
}
