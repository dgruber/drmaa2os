package libdrmaa

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dgruber/drmaa2interface"
)

func QstatJobState(jobID string) (string, error) {
	out, err := exec.Command("qstat", "-u", "*").Output()
	if err != nil {
		return "", fmt.Errorf("error executing qstat -u *: %v", err)
	}
	stateMap := ParseQstatForJobIDs(string(out), []string{jobID})
	state, exists := stateMap[jobID]
	if !exists {
		return "", fmt.Errorf("does not exist")
	}
	return state, nil
}

func QstatGetJobIDs() ([]string, error) {
	out, err := exec.Command("qstat", "-u", "*").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing qstat -u *: %v", err)
	}
	return ParseQstatForAllJobIDs(string(out)), nil
}

func ParseQstatForAllJobIDs(out string) []string {
	res := make([]string, 0, 16)

	// Example:
	//
	//# qstat -u \*
	// job-ID  prior   name       user         state submit/start at     queue                          slots ja-task-ID
	// -----------------------------------------------------------------------------------------------------------------
	// 900 0.55500 sleep      root         r     01/22/2022 16:46:04 all.q@master                       1
	// 901 0.55500 sleep      root         r     01/22/2022 16:46:05 all.q@master                       1

	// skip first two header lines
	lines := strings.Split(out, "\n")
	for i := 2; i < len(lines); i++ {
		// remove heading spaces and split
		values := strings.Fields(lines[i])
		if len(values) > 1 && values[0] != "" {
			res = append(res, values[0])
		}
	}
	return res
}

func ParseQstatForJobIDs(out string, ids []string) map[string]string {
	hm := make(map[string]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		hm[ids[i]] = nil
	}
	res := map[string]string{}
	// Example:
	//# qstat -u \*
	// job-ID  prior   name       user         state submit/start at     queue                          slots ja-task-ID
	// -----------------------------------------------------------------------------------------------------------------
	// 900 0.55500 sleep      root         r     01/22/2022 16:46:04 all.q@master                       1
	// 901 0.55500 sleep      root         r     01/22/2022 16:46:05 all.q@master                       1

	// skip first two header lines
	lines := strings.Split(out, "\n")
	for i := 2; i < len(lines); i++ {
		// remove heading spaces and split
		values := strings.Fields(lines[i])
		if len(values) > 1 && values[0] != "" {
			if _, exists := hm[values[0]]; exists {
				// we are interested in the job
				if len(values) > 5 {
					// jobID = state
					res[values[0]] = values[4]
				}
			}
		}
	}
	return res
}

func ConvertQstatJobState(state string) drmaa2interface.JobState {
	switch state {
	case "r":
		return drmaa2interface.Running
	case "t":
		return drmaa2interface.Running
	case "R":
		return drmaa2interface.Running
	case "hr":
		return drmaa2interface.Running
	case "h":
		return drmaa2interface.QueuedHeld
	case "w":
		return drmaa2interface.Queued
	case "P":
		return drmaa2interface.Queued
	case "N":
		return drmaa2interface.Queued
	case "E":
		return drmaa2interface.Failed
	case "Eqw":
		return drmaa2interface.Failed
	case "qw":
		return drmaa2interface.Queued
	case "Hqw":
		return drmaa2interface.QueuedHeld
	case "Tr":
		return drmaa2interface.QueuedHeld
	case "T":
		return drmaa2interface.QueuedHeld
	case "S":
		return drmaa2interface.Suspended
	case "s":
		return drmaa2interface.Suspended
	}
	return drmaa2interface.Undetermined
}
