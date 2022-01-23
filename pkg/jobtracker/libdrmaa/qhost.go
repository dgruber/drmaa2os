package libdrmaa

import (
	"fmt"
	"os/exec"
	"strings"
)

func QhostGetAllHosts() ([]string, error) {
	out, err := exec.Command("qhost").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing qhost: %v", err)
	}
	return ParseQhostForHostnames(string(out)), nil
}

func ParseQhostForHostnames(out string) []string {
	// Example:
	//# qhost
	//	HOSTNAME                ARCH         NCPU NSOC NCOR NTHR  LOAD  MEMTOT  MEMUSE  SWAPTO  SWAPUS
	//----------------------------------------------------------------------------------------------
	//global                  -               -    -    -    -     -       -       -       -       -
	//master                  lx-amd64        4    1    4    4  0.07    1.9G  442.2M 1024.0M     0.0                    1

	// skip first two header lines and "global" host
	lines := strings.Split(out, "\n")
	hostnames := make([]string, 0, len(lines))
	for i := 3; i < len(lines); i++ {
		// remove heading spaces and split
		values := strings.Split(strings.Trim(lines[i], " "), " ")
		if len(values) > 1 && values[0] != "" {
			hostnames = append(hostnames, values[0])
		}
	}
	return hostnames
}
