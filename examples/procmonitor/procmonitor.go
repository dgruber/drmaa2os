package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

func main() {

	var sm drmaa2interface.SessionManager
	var err error

	if os.Getenv("DOCKER") == "TRUE" {
		sm, err = drmaa2os.NewDockerSessionManager("testdocker.db")
		if err != nil {
			log.Panic(err)
		}
	} else if os.Getenv("KUBERNETES") == "TRUE" {
		sm, err = drmaa2os.NewKubernetesSessionManager(nil, "testkubernetes.db")
		if err != nil {
			log.Panic(err)
		}
	} else {
		sm, err = drmaa2os.NewDefaultSessionManager("testprocess.db")
		if err != nil {
			log.Panic(err)
		}
	}

	monitor, err := sm.OpenMonitoringSession("monitor")
	if err != nil {
		log.Panic(err)
	}
	defer monitor.CloseMonitoringSession()

	host, err := monitor.GetAllMachines(nil)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Host list:")
	for i := range host {
		fmt.Printf("Host name: %v\n", host[i].Name)
		fmt.Printf("Load (1 min avg): %v\n", host[i].Load)
		fmt.Printf("Physical memory (in mb): %d\n", host[i].PhysicalMemory/1024/1024)
		fmt.Printf("Virtual memory (in mb): %d\n", host[i].VirtualMemory/1024/1024)
		fmt.Printf("OS: %s\n", host[i].OS.String())
		fmt.Printf("OS version: %s\n", host[i].OSVersion.String())
		fmt.Printf("Host extensions:\n")
		for k, v := range host[i].ExtensionList {
			fmt.Printf("%s: %s\n", k, v)
		}
	}

	jobs, err := monitor.GetAllJobs(drmaa2interface.CreateJobInfo())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("found %d jobs\n", len(jobs))
	for _, job := range jobs {
		fmt.Printf("job id: %s, state: %s\n", job.GetID(), job.GetState().String())
	}

	if len(jobs) > 0 {
		fmt.Printf("Detailed information about first job:\n")
		jobInfo, err := jobs[0].GetJobInfo()
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("JobID: %s\n", jobInfo.ID)
		fmt.Printf("Owner: %s\n", jobInfo.JobOwner)
		fmt.Printf("State: %s\n", jobInfo.State.String())
		fmt.Printf("Runtime (wallclock): %s\n", jobInfo.WallclockTime.String())

		fmt.Printf("JobInfo extensions:\n")
		for k, v := range jobInfo.ExtensionList {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
}
