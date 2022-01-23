package libdrmaa

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dgruber/drmaa/gestatus"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/d2hlp"
)

// Implements the monitorer interface so that a monitoring session
// can be generated. As drmaa (v1) does not support that, the command
// line utilities (qstat) are used for getting the neccessary information.
// That has performance and compatibility impacts as different systems
// use different tooling. First support is implemented for open source
// Grid Engine.

func (m *DRMAATracker) OpenMonitoringSession(name string) error {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return fmt.Errorf("unsupported workload manager")
	}
	return nil
}

func (m *DRMAATracker) CloseMonitoringSession(name string) error {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return fmt.Errorf("unsupported workload manager")
	}
	return nil
}

func (m *DRMAATracker) GetAllJobIDs(filter *drmaa2interface.JobInfo) ([]string, error) {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return nil, fmt.Errorf("unsupported workload manager")
	}
	jobids, err := QstatGetJobIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to get job IDs from qstat: %v", err)
	}
	// TODO filter
	if filter != nil {
		return jobids, fmt.Errorf("filter not yet implemented in GetAllJobIDs")
	}
	return jobids, nil
}

func (m *DRMAATracker) GetAllQueueNames(names []string) ([]string, error) {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return nil, fmt.Errorf("unsupported workload manager")
	}
	queueList, err := QconfSQL()
	if err != nil {
		return nil, err
	}
	if names == nil {
		return queueList, nil
	}
	filter := d2hlp.NewStringFilter(queueList)
	return filter.GetIncludedSubset(names), nil
}

func (m *DRMAATracker) GetAllMachines(names []string) ([]drmaa2interface.Machine, error) {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return nil, fmt.Errorf("unsupported workload manager")
	}
	// TODO parse host details
	machines, err := QhostGetAllHosts()
	if err != nil {
		return nil, err
	}
	// no filter
	if names == nil {
		return d2hlp.ConvertStringsToMachines(machines), nil
	}
	// filter
	filter := d2hlp.NewStringFilter(machines)
	return d2hlp.ConvertStringsToMachines(filter.GetIncludedSubset(names)), nil
}

func (m *DRMAATracker) JobInfoFromMonitor(id string) (drmaa2interface.JobInfo, error) {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return drmaa2interface.JobInfo{}, fmt.Errorf("unsupported workload manager")
	}
	jobStatus, err := gestatus.GetJob(id)
	if err != nil {
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed to get job status with qstat -xml: %v", err)
	}
	ji := drmaa2interface.JobInfo{}
	ji.ID = id
	ji.JobOwner = jobStatus.JobOwner()
	ji.AllocatedMachines = jobStatus.DestinationHostList()
	ji.WallclockTime = jobStatus.RunTime()
	ji.DispatchTime = jobStatus.StartTime()
	ji.SubmissionTime = jobStatus.SubmissionTime()
	ji.ExtensionList = map[string]string{
		"account":  jobStatus.JobAccountName(),
		"jobclass": jobStatus.JobClassName(),
		"mail":     strings.Join(jobStatus.MailAdresses(), ","),
	}
	for i := 0; i < len(jobStatus.DestinationSlotsList()); i++ {
		s, _ := strconv.Atoi(jobStatus.DestinationSlotsList()[i])
		ji.Slots += int64(s)
	}

	if len(jobStatus.DestinationQueueInstanceList()) >= 1 {
		if qi := strings.Split(jobStatus.DestinationQueueInstanceList()[0], "@"); len(qi) == 2 {
			ji.QueueName = strings.Split(jobStatus.DestinationQueueInstanceList()[0], "@")[1]
		}
	}
	// TODO add remaining values
	return ji, nil
}
