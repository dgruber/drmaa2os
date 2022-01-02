package simpletracker

import (
	"fmt"
	"os"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/shirou/gopsutil/v3/process"
)

func GetAllProcesses() ([]string, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	jobIDs := make([]string, 0, len(processes))
	for _, proc := range processes {
		jobIDs = append(jobIDs, fmt.Sprintf("%d", proc.Pid))
	}
	return jobIDs, nil
}

func GetJobInfo(id int32) (drmaa2interface.JobInfo, error) {
	exists, err := process.PidExists(id)
	if err != nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed to check PIDs existence: %v", err)
	}
	if exists == false {
		return drmaa2interface.JobInfo{
			State:    drmaa2interface.Undetermined,
			SubState: "process not found",
		}, fmt.Errorf("process not found")
	}
	proc, err := process.NewProcess(id)
	if err != nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed to get internal process: %v", err)
	}
	return ProcessToJobInfo(proc), nil
}

func ProcessToJobInfo(proc *process.Process) drmaa2interface.JobInfo {
	var ji drmaa2interface.JobInfo

	startedMsSinceEpochproc, _ := proc.CreateTime()
	ji.DispatchTime = time.Unix(0, startedMsSinceEpochproc*int64(time.Millisecond))
	ji.SubmissionTime = ji.DispatchTime

	ji.JobOwner, _ = proc.Username()

	ji.WallclockTime = time.Now().Sub(ji.DispatchTime)

	ji.State = drmaa2interface.Running
	ji.ID = fmt.Sprintf("%d", proc.Pid)
	hostname, _ := os.Hostname()
	ji.AllocatedMachines = []string{hostname}

	var affinity string
	cpuaffinity, err := proc.CPUAffinity()
	if err == nil {
		for _, i := range cpuaffinity {
			affinity = fmt.Sprintf("%s%d ", affinity, i)
		}
	}

	// a few extensions
	extensions := map[string]string{}

	extensions["name"], _ = proc.Name()

	if cli, err := proc.Cmdline(); err == nil && cli != "" {
		extensions[jobtracker.DRMAA2_MS_JOBINFO_COMMANDLINE] = cli
	}

	if workdir, err := proc.Cwd(); err == nil && workdir != "" {
		extensions[jobtracker.DRMAA2_MS_JOBINFO_WORKINGDIR] = workdir
	}

	usage, _ := proc.CPUPercent()
	extensions["cpu_usage"] = fmt.Sprintf("%f", usage)

	if affinity != "" {
		extensions["cpu_affinity"] = affinity
	}

	mem, err := proc.MemoryInfo()
	if err == nil {
		extensions["memory_usage"] = mem.String()
		extensions["memory_usage_rss"] = fmt.Sprintf("%d", mem.RSS)
		extensions["memory_usage_vms"] = fmt.Sprintf("%d", mem.VMS)
	}

	ji.ExtensionList = extensions

	return ji
}
