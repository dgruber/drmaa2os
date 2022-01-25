package libdrmaa

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dgruber/drmaa2interface"
)

func QhostGetAllHosts() ([]drmaa2interface.Machine, error) {
	out, err := exec.Command("qhost").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing qhost: %v", err)
	}
	return ParseQhostForHostnames(string(out)), nil
}

func ParseQhostForHostnames(out string) []drmaa2interface.Machine {
	// Example:
	//# qhost
	//	HOSTNAME                ARCH         NCPU NSOC NCOR NTHR  LOAD  MEMTOT  MEMUSE  SWAPTO  SWAPUS
	//----------------------------------------------------------------------------------------------
	//global                  -               -    -    -    -     -       -       -       -       -
	//master                  lx-amd64        4    1    4    4  0.07    1.9G  442.2M 1024.0M     0.0                    1

	// skip first two header lines and "global" host
	lines := strings.Split(out, "\n")
	machines := make([]drmaa2interface.Machine, 0, len(lines))
	for i := 3; i < len(lines); i++ {
		// remove heading spaces and split
		values := strings.Fields(strings.Trim(lines[i], " "))
		if len(values) > 1 && values[0] != "" {
			machine := drmaa2interface.Machine{
				Name: values[0],
			}
			// if machine has no load it is not available
			if values[6] == "-" {
				machine.Available = false
			} else {
				machine.Available = true
				machine.Load, _ = strconv.ParseFloat(values[6], 64)
			}
			if values[1] == "lx-amd64" {
				machine.Architecture = drmaa2interface.IA64
				machine.OS = drmaa2interface.Linux
			} else {
				machine.Architecture = drmaa2interface.OtherCPU
			}
			machine.Sockets, _ = strconv.ParseInt(values[3], 10, 64)
			if machine.Sockets <= 0 {
				machine.Sockets = 1
			}
			cores, _ := strconv.ParseInt(values[4], 10, 64)
			if cores <= 0 {
				cores = 1
			}
			threads, _ := strconv.ParseInt(values[5], 10, 64)
			if threads <= 0 {
				threads = 1
			}
			machine.CoresPerSocket = cores / machine.Sockets
			machine.ThreadsPerCore = threads / cores
			// TODO calculate the memory - must be already in some other project
			machines = append(machines, machine)
		}
	}
	return machines
}
