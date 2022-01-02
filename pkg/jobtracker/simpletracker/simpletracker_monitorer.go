package simpletracker

import (
	"fmt"
	"strconv"

	"github.com/dgruber/drmaa2interface"
)

func (m *JobTracker) OpenMonitoringSession(name string) error {
	return nil
}

func (m *JobTracker) CloseMonitoringSession(name string) error {
	return nil
}

func (m *JobTracker) GetAllJobIDs(filter *drmaa2interface.JobInfo) ([]string, error) {
	processList, err := GetAllProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to get job IDs: %v", err)
	}
	// TODO filter
	if filter != nil {
		return processList, fmt.Errorf("filter not yet implemented in GetAllJobIDs")
	}
	return processList, nil
}

func (m *JobTracker) GetAllQueueNames(names []string) ([]string, error) {
	return []string{}, nil
}

func (m *JobTracker) GetAllMachines(names []string) ([]drmaa2interface.Machine, error) {
	localMachine, err := GetLocalMachineInfo()
	if err != nil {
		return nil, err
	}
	// no filter
	if names == nil {
		return []drmaa2interface.Machine{localMachine}, nil
	}
	for _, allowedName := range names {
		if localMachine.Name != allowedName {
			// filter
			continue
		} else {
			return []drmaa2interface.Machine{localMachine}, nil
		}
	}
	// filter does not include localhost name
	return []drmaa2interface.Machine{}, nil
}

func (m *JobTracker) JobInfoFromMonitor(id string) (drmaa2interface.JobInfo, error) {
	proccessID, err := strconv.Atoi(id)
	if err != nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("job ID is not a valid process ID: %v", err)
	}
	return GetJobInfo(int32(proccessID))
}
