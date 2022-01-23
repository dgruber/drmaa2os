package libdrmaa

import (
	"fmt"
	"os/exec"
	"strings"
)

func QstatGetJobIDs() ([]string, error) {
	out, err := exec.Command("qstat", "-u", "*").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing qstat -u *: %v", err)
	}
	return ParseQstatForJobIDs(string(out)), nil
}

func ParseQstatForJobIDs(out string) []string {
	// Example:
	//# qstat -u \*
	// job-ID  prior   name       user         state submit/start at     queue                          slots ja-task-ID
	// -----------------------------------------------------------------------------------------------------------------
	// 900 0.55500 sleep      root         r     01/22/2022 16:46:04 all.q@master                       1
	// 901 0.55500 sleep      root         r     01/22/2022 16:46:05 all.q@master                       1

	// skip first two header lines
	lines := strings.Split(out, "\n")
	jobIDs := make([]string, 0, len(lines))
	for i := 2; i < len(lines); i++ {
		// remove heading spaces and split
		values := strings.Split(strings.Trim(lines[i], " "), " ")

		if len(values) > 1 && values[0] != "" {
			jobIDs = append(jobIDs, values[0])
		}
	}
	return jobIDs
}
