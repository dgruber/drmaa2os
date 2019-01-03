package main

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
)

func CreateOpenJobSession(sm drmaa2interface.SessionManager, name, contact string) (drmaa2interface.JobSession, error) {
	js, err := sm.CreateJobSession(name, contact)
	if err != nil {
		return sm.OpenJobSession(name)
	}
	return js, err
}

func main() {
	sm, err := drmaa2os.NewDefaultSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := CreateOpenJobSession(sm, "jobsession1", "")
	if err != nil {
		panic(err)
	}
	defer js.Close()
	defer sm.DestroyJobSession("jobsession1")

	jt := drmaa2interface.JobTemplate{
		JobName:       "job1",
		RemoteCommand: "sleep",
		Args:          []string{"1"},
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)
	jobinfo, err := job.GetJobInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", jobinfo)
}
