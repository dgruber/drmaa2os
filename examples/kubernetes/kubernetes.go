package main

import (
	"encoding/base64"
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
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

	// Allocate a SessionManager managing jobs in the "default"
	// namespace, with no pre-initialized ClientSet (so that a
	// new one gets allocated). The local DB is used by the
	// SessionManager to make some job details persistent.
	sm, err := drmaa2os.NewKubernetesSessionManager(
		kubernetestracker.KubernetesTrackerParameters{
			Namespace: "default",
			ClientSet: nil,
		}, "testdb.db")
	if err != nil {
		panic(err)
	}

	js := createJobSession(sm)
	defer js.Close()

	jt := drmaa2interface.JobTemplate{
		// JobName must be unique or not set ("").
		RemoteCommand: "/bin/sh",
		JobCategory:   "busybox:latest",
		Args:          []string{"-c", `cp /input/data.txt /persistent/data.txt`},
	}

	// Use storage from existing PVC "nfs-pvc" and mount inside the container to
	// directory "/persistent". Note that the storage needs to support multi-write
	// as multiple jobs are accessing it concurrently here. Examples of this
	// kind of storage are Google Filestore (for GKE) or a self-provisioned
	// NFS server.
	//
	// Stage data into job which gets mounted under /input/data.txt
	// as a ConfigMap
	jt.StageInFiles = map[string]string{
		"/persistent": "pvc:nfs-pvc",
		"/input/data.txt": "configmap-data:" +
			base64.StdEncoding.EncodeToString([]byte("my input data set")),
	}

	fmt.Println("running data pre-processing job")

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

	// removing job and config map object from Kubernetes
	job.Reap()

	fmt.Println("Starting job array with 17 jobs each transforming on char to upper case")
	jt.StageInFiles = map[string]string{
		"/persistent": "pvc:nfs-pvc",
	}
	jt.Args = []string{"-c", `cut -c$(TASK_ID) /persistent/data.txt | tr [:lower:] [:upper:] > /persistent/data_$(TASK_ID).txt`}

	jobs, err := js.RunBulkJobs(jt, 1, 17, 1, 17)
	if err != nil {
		panic(err)
	}
	for _, j := range jobs.GetJobs() {
		j.WaitTerminated(drmaa2interface.InfiniteTime)
		fmt.Printf("Job %s finished\n", j.GetID())
		j.Reap()
	}

	fmt.Println("All data processing jobs finished. Starting post-processing.")

	// output data is not in order! :)
	jt.Args = []string{"-c", `cat /persistent/data_*.txt > /persistent/output_data.txt`}
	job, err = js.RunJob(jt)
	if err != nil {
		panic(err)
	}
	job.WaitTerminated(drmaa2interface.InfiniteTime)
	fmt.Printf("Post-processing job finished in state: %s\n", job.GetState().String())
	job.Reap()
}
