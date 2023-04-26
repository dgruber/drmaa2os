package main

import (
	"encoding/json"
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
)

func main() {

	sm, err := drmaa2os.NewKubernetesSessionManager(
		kubernetestracker.KubernetesTrackerParameters{
			Namespace: "default",
		}, "testdb.db")
	if err != nil {
		panic(err)
	}

	js := createJobSession(sm)
	defer js.Close()

	jt := drmaa2interface.JobTemplate{
		JobCategory:   "nvcr.io/nvidia/pytorch:23.04-py3",
		RemoteCommand: "nvidia-smi",
		OutputPath:    "/tmp/output.txt",
		ErrorPath:     "/tmp/output.txt",
		Extension: drmaa2interface.Extension{
			ExtensionList: map[string]string{
				"distribution": "gke",                 // "aks", "eks", "gke"
				"accelerator":  "1*nvidia-tesla-v100", // amount*type
				"pullpolicy":   "Always",              // "IfNotPresent", "Never"
			},
		},
	}

	fmt.Println("Running pytorch job")

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job submitted successfully. Waiting until the job is finished.")

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	ji, err := job.GetJobInfo()
	if err != nil {
		panic(err)
	}

	if ji.ExtensionList != nil && ji.ExtensionList["output"] != "" {
		fmt.Println("Job Output:")
		fmt.Println(ji.ExtensionList["output"])
	}

	fmt.Println("Job Info:")
	jobInfo, _ := json.Marshal(ji)
	fmt.Println(string(jobInfo))

	// removing job object and config map object from Kubernetes
	job.Reap()
}

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
