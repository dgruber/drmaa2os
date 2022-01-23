package d2hlp

import (
	"time"

	"github.com/dgruber/drmaa2interface"
)

// JobInfoIsUnset returns true when the job info struct does not
// filter any jobs, i.e. all fields are set to the specified unset
// value. Un Unset JobInfo struct is returned by drmaa2interface.CreateJobInfo().
func JobInfoIsUnset(ji drmaa2interface.JobInfo) bool {
	if ji.ID != "" {
		return false
	}
	if ji.ExitStatus != drmaa2interface.UnsetNum {
		return false
	}
	if ji.TerminatingSignal != "" {
		return false
	}
	if ji.Annotation != "" {
		return false
	}
	if ji.State != drmaa2interface.Unset {
		return false
	}
	if ji.SubState != "" {
		return false
	}
	if ji.AllocatedMachines != nil {
		return false
	}
	if ji.SubmissionMachine != "" {
		return false
	}
	if ji.JobOwner != "" {
		return false
	}
	if ji.Slots != drmaa2interface.UnsetNum {
		return false
	}
	if ji.QueueName != "" {
		return false
	}
	if ji.WallclockTime != 0 {
		return false
	}
	if ji.CPUTime != drmaa2interface.UnsetTime {
		return false
	}
	var nullTime time.Time
	if ji.SubmissionTime != nullTime {
		return false
	}
	if ji.DispatchTime != nullTime {
		return false
	}
	if ji.FinishTime != nullTime {
		return false
	}
	return true
}

// JobInfoMatches returns true when the given job info is allowed
// by the filter.
func JobInfoMatches(ji drmaa2interface.JobInfo, filter drmaa2interface.JobInfo) bool {
	if filter.ID != "" {
		if ji.ID != filter.ID {
			return false
		}
	}
	if filter.ExitStatus != drmaa2interface.UnsetNum {
		if ji.ExitStatus != filter.ExitStatus {
			return false
		}
	}
	if filter.TerminatingSignal != "" {
		if ji.TerminatingSignal != filter.TerminatingSignal {
			return false
		}
	}
	if filter.Annotation != "" {
		if ji.Annotation != filter.Annotation {
			return false
		}
	}
	if filter.State != drmaa2interface.Unset {
		if ji.State != filter.State {
			return false
		}
	}
	if filter.SubState != "" {
		if ji.SubState != filter.SubState {
			return false
		}
	}
	if filter.AllocatedMachines != nil {
		// must run on a superset of the given machines
		if len(ji.AllocatedMachines) < len(filter.AllocatedMachines) {
			return false
		}

		for _, v := range filter.AllocatedMachines {
			found := false
			for _, i := range ji.AllocatedMachines {
				if v == i {
					found = true
					break
				}
			}
			if found == false {
				return false
			}
		}
	}
	if filter.SubmissionMachine != "" {
		if ji.SubmissionMachine != filter.SubmissionMachine {
			return false
		}
	}
	if filter.JobOwner != "" {
		if ji.JobOwner != filter.JobOwner {
			return false
		}
	}
	if filter.Slots != drmaa2interface.UnsetNum {
		if ji.Slots != filter.Slots {
			return false
		}
	}
	if filter.QueueName != "" {
		if ji.QueueName != filter.QueueName {
			return false
		}
	}
	if filter.WallclockTime != 0 {
		if ji.WallclockTime < filter.WallclockTime {
			return false
		}
	}
	if filter.CPUTime != drmaa2interface.UnsetTime {
		if ji.CPUTime < filter.CPUTime {
			return false
		}
	}
	var nullTime time.Time
	if filter.SubmissionTime != nullTime {
		if ji.SubmissionTime.Before(filter.SubmissionTime) {
			return false
		}
	}
	if filter.DispatchTime != nullTime {
		if ji.DispatchTime.Before(filter.DispatchTime) {
			return false
		}
	}
	if filter.FinishTime != nullTime {
		if ji.FinishTime.Before(filter.FinishTime) {
			return false
		}
	}
	return true
}

// ConvertStringsToMachines converts machine names into drmaa2 machines
// in which only the name is set.
func ConvertStringsToMachines(names []string) []drmaa2interface.Machine {
	res := make([]drmaa2interface.Machine, 0, len(names))
	for i := 0; i < len(names); i++ {
		res = append(res, drmaa2interface.Machine{
			Name: names[i],
		})
	}
	return res
}

// StringFilter implements a lookup method for strings
type StringFilter struct {
	values map[string]interface{}
}

// NewStringFilter creates a hashmap for efficiently looking
// up if a value is included in the map or return a subset
// of a given filter.
func NewStringFilter(values []string) *StringFilter {
	sf := StringFilter{
		values: make(map[string]interface{}, len(values)),
	}
	for i := range values {
		sf.values[values[i]] = nil
	}
	return &sf
}

// IsIncluded returns true when the item is found in the filter list.
func (sf *StringFilter) IsIncluded(filter string) bool {
	_, exists := sf.values[filter]
	return exists
}

func (sf *StringFilter) GetIncludedSubset(filter []string) []string {
	result := make([]string, 0, len(filter))
	for i := 0; i < len(filter); i++ {
		if sf.IsIncluded(filter[i]) {
			result = append(result, filter[i])
		}
	}
	return result
}
