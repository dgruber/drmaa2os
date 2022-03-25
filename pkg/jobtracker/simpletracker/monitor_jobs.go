package simpletracker

import (
	"fmt"
	"os"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/extension"
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

	// a few extensions
	extensions := map[string]string{}

	extensions[extension.JobInfoDefaultMSessionProcessName], _ = proc.Name()

	if cli, err := proc.Cmdline(); err == nil && cli != "" {
		extensions[extension.JobInfoDefaultMSessionCommandLine] = cli
	}

	if workdir, err := proc.Cwd(); err == nil && workdir != "" {
		extensions[extension.JobInfoDefaultMSessionWorkingDir] = workdir
	}

	usage, _ := proc.CPUPercent()
	extensions[extension.JobInfoDefaultMSessionCPUUsage] = fmt.Sprintf("%f", usage)

	var affinity string
	cpuaffinity, err := proc.CPUAffinity()
	if err == nil {
		for _, i := range cpuaffinity {
			affinity = fmt.Sprintf("%s%d ", affinity, i)
		}
	}

	if affinity != "" {
		extensions[extension.JobInfoDefaultMSessionCPUAffinity] = affinity
	}

	mem, err := proc.MemoryInfo()
	if err == nil {
		extensions[extension.JobInfoDefaultMSessionMemoryUsage] = mem.String()
		extensions[extension.JobInfoDefaultMSessionMemoryUsageRSS] = fmt.Sprintf("%d", mem.RSS)
		extensions[extension.JobInfoDefaultMSessionMemoryUsageVMS] = fmt.Sprintf("%d", mem.VMS)
	}

	ji.ExtensionList = extensions

	return ji
}
