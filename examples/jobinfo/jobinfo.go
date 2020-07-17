package main

import (
	"fmt"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
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

	fmt.Printf("ID: %s\n", jobinfo.ID)
	fmt.Printf("State: %s\n", jobinfo.State)
	fmt.Printf("SubState: %s\n", jobinfo.SubState)
	fmt.Printf("Annotation: %s\n", jobinfo.Annotation)
	fmt.Printf("ExitStatus: %d\n", jobinfo.ExitStatus)
	fmt.Printf("TerminatingSignal: %s\n", jobinfo.TerminatingSignal)
	fmt.Printf("AllocatedMachines: %v\n", jobinfo.AllocatedMachines)
	fmt.Printf("SubmissionMachine: %s\n", jobinfo.SubmissionMachine)
	fmt.Printf("JobOwner: %s\n", jobinfo.JobOwner)
	fmt.Printf("Slots: %d\n", jobinfo.Slots)
	fmt.Printf("QueueName: %s\n", jobinfo.QueueName)
	fmt.Printf("WallclockTime: %s\n", jobinfo.WallclockTime)
	fmt.Printf("CPUTime: %d\n", jobinfo.CPUTime)
	fmt.Printf("SubmissionTime: %s\n", jobinfo.SubmissionTime)
	fmt.Printf("DispatchTime: %s\n", jobinfo.DispatchTime)
	fmt.Printf("FinishTime: %s\n", jobinfo.FinishTime)
}
