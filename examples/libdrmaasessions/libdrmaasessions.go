package main

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
)

func main() {
	params := libdrmaa.LibDRMAASessionParams{
		ContactString:           "",
		UsePersistentJobStorage: true,
		DBFilePath:              "testdbjobs.db",
	}
	sm, err := drmaa2os.NewLibDRMAASessionManagerWithParams(params, "testdb.db")
	if err != nil {
		panic(err)
	}

	var contact string
	if len(os.Args) == 2 {
		contact = os.Args[1]
		fmt.Printf("using contact string %s\n", contact)
	}

	js, err := sm.CreateJobSession("jobsession6", contact)
	if err != nil {
		fmt.Printf("failed creating job session, trying to open it\n")
		js, err = sm.OpenJobSession("jobsession6")
		if err != nil {
			panic(err)
		}
	} else {
		contact, err := js.GetContact()
		if err != nil {
			panic(err)
		}
		fmt.Printf("session has contact string %s\n", contact)
	}
	defer js.Close()

	jobinfo := drmaa2interface.CreateJobInfo()
	jobs, err := js.GetJobs(jobinfo)
	if err != nil {
		panic(err)
	}
	for _, job := range jobs {
		fmt.Printf("found job %s\n", job.GetID())
	}

	job, err := js.RunJob(drmaa2interface.JobTemplate{
		JobName:       "job1",
		RemoteCommand: "/bin/sleep",
		Args:          []string{"100"},
	})
	if err != nil {
		fmt.Printf("job submission failed: %v\n", err)
		js.Close()
		os.Exit(1)
	}
	fmt.Printf("job submitted with ID %s\n", job.GetID())
}
