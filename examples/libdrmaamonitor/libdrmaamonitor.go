package main

import (
	"fmt"
	"time"

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
	// in Grid Engine drmaav1 either one MonitoringSession or one JobSession can
	// be used but not multiple ones...
	ms, err := sm.OpenMonitoringSession("monitoringsession")
	if err != nil {
		panic(err)
	}
	defer ms.CloseMonitoringSession()

	for range time.NewTicker(time.Second * 1).C {
		machines, err := ms.GetAllMachines(nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("found machines: %v\n", machines)

		queues, err := ms.GetAllQueues(nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("found queues: ")
		for _, queue := range queues {
			fmt.Printf("%s ", queue.Name)
		}
		fmt.Println()

		jobs, err := ms.GetAllJobs(drmaa2interface.CreateJobInfo())
		if err != nil {
			panic(err)
		}
		for _, job := range jobs {
			fmt.Printf("monitor: found job %s in state %s\n",
				job.GetID(), job.GetState().String())
		}
	}
}
