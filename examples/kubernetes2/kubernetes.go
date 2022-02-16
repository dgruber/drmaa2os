package main

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/extension"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
)

func main() {
	sm := createSessionManager()
	js := createJobSession(sm)
	defer js.Close()

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "/bin/sh",
		JobCategory:   "busybox:latest",
		JobEnvironment: map[string]string{
			"env": "var",
		},
		Args:       []string{"-c", `env && ls /container/dir/ && touch /container/dir/output.txt`},
		OutputPath: "/dev/stdout",
		StageInFiles: map[string]string{
			// mount host temp directory into container under /container/dir
			"/container/dir": extension.JobTemplateK8sStageInFromHostPathPrefix + "/tmp",
		},
	}

	jt.ExtensionList = map[string]string{
		// populate additional envs from a secret "secret-test" - which
		// must exist before:
		//extension.JobTemplateK8sEnvFromSecret: "secret-test",
		// run in priviledged mode (default is false) to allow
		// to access the host filesystem (StageIn/StageOut)
		extension.JobTemplateK8sPrivileged: "TRUE",
	}

	fmt.Println("running job")

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
	fmt.Printf("job output: %s\n", ji.ExtensionList[extension.JobInfoK8sJSessionJobOutput])

	// removing created artifacts from Kubernetes (including job object)
	//job.Reap()
}

func createSessionManager() drmaa2interface.SessionManager {
	sm, err := drmaa2os.NewKubernetesSessionManager(
		kubernetestracker.KubernetesTrackerParameters{
			Namespace: "default",
			ClientSet: nil,
		}, "testdb.db")
	if err != nil {
		panic(err)
	}
	return sm
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

func print(ji drmaa2interface.JobInfo) {
	fmt.Printf("Submission time: %s\n", ji.SubmissionMachine)
	fmt.Printf("Dispatch time: %s\n", ji.DispatchTime)
	fmt.Printf("End time: %s\n", ji.FinishTime)
	fmt.Printf("State: %s\n", ji.State)
	fmt.Printf("Job ID: %s\n", ji.ID)
}
