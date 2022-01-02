package simpletracker

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgruber/drmaa2interface"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// GetLocalMachineInfo collects information about the local machine
// and returns a current DRMAA2 machine info struct.
func GetLocalMachineInfo() (drmaa2interface.Machine, error) {
	machine := drmaa2interface.Machine{}

	machine.Name, _ = os.Hostname()
	machine.Available = true

	mem, _ := mem.VirtualMemory()
	if mem != nil {
		machine = AddMemory(machine, mem)
	}

	// "This attributes describes the 1-minute average load on the given machine"
	avgStat, _ := load.Avg()
	machine.Load = avgStat.Load1

	hostInfo, _ := host.Info()
	if hostInfo != nil {
		machine = AddHostInfo(machine, hostInfo)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return machine, err
	}
	machine.Sockets, machine.CoresPerSocket, machine.ThreadsPerCore, err = CollectSocketCoreThreads(cpuInfo)

	// a few extensions
	uptime, _ := host.Uptime()

	machine.ExtensionList = map[string]string{
		"tempdir":        os.TempDir(),
		"load5min":       fmt.Sprintf("%f", avgStat.Load5),
		"load15min":      fmt.Sprintf("%f", avgStat.Load15),
		"uptime_seconds": fmt.Sprintf("%d", uptime),
	}

	return machine, err
}

func AddMemory(machine drmaa2interface.Machine, mem *mem.VirtualMemoryStat) drmaa2interface.Machine {
	// "This attribute describes the amount of virtual memory in kilobyte
	// available for a job executing on this machine"
	machine.VirtualMemory = int64(mem.Total + mem.SwapTotal)
	// "This attribute describes the amount of physical memory in kilobyte
	// installed in this machine."
	machine.PhysicalMemory = int64(mem.Total)
	return machine
}

func AddHostInfo(machine drmaa2interface.Machine, hostInfo *host.InfoStat) drmaa2interface.Machine {
	if hostInfo.OS == "linux" {
		machine.OS = drmaa2interface.Linux
	} else if hostInfo.OS == "darwin" {
		machine.OS = drmaa2interface.MacOS
	} else if hostInfo.OS == "freebsd" {
		machine.OS = drmaa2interface.BSD
	}

	// TODO other archs?
	machine.Architecture = drmaa2interface.X64

	if machine.OS == drmaa2interface.MacOS {
		// on darwin:
		// kernel version: 21.2.0
		// platform version: 12.1
		version := strings.Split(hostInfo.PlatformVersion, ".")
		if len(version) == 2 {
			machine.OSVersion = drmaa2interface.Version{
				Major: version[0],
				Minor: version[1],
			}
		}
	} else {
		// TODO needs further tests
		version := strings.Split(hostInfo.PlatformVersion, ".")
		if len(version) == 2 {
			machine.OSVersion = drmaa2interface.Version{
				Major: version[0],
				Minor: version[1],
			}
		}
	}
	return machine
}

// CollectSocketCoreThreads returns the amount of sockets, cores per socket,
// and threads per core.
func CollectSocketCoreThreads(cpuInfo []cpu.InfoStat) (int64, int64, int64, error) {
	sockets := int32(0)
	firstCPU := int32(9999)
	coresPerSocket := 0

	physicalIDs := make(map[string]interface{})
	coreIDs := make(map[string]interface{})

	for _, info := range cpuInfo {
		// is first socket
		if firstCPU == 9999 {
			firstCPU = info.CPU
		}
		if firstCPU == info.CPU {
			// count cores of the first CPU ID
			coresPerSocket++
			physicalIDs[info.PhysicalID] = nil
			coreIDs[info.CoreID] = nil
		}
		if (info.CPU + 1) > sockets {
			sockets = (info.CPU + 1)
		}
	}

	threadsPerCore := len(coreIDs) / len(physicalIDs)

	return int64(sockets), int64(len(physicalIDs)), int64(threadsPerCore), nil
}
