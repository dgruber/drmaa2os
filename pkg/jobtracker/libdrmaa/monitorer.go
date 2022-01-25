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
// Grid Engine. Note, that in GE only one job session or one monitoring
// session can be used at one point in time. This is a drmaa (v1)
// limitation...

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

// GetAllMachines returns all machines the cluster consists of.
// If names is != nil, it returns only a subset of machines which
// names are defined in names and are in the cluster.
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
	// filter
	if names != nil {
		res := make([]drmaa2interface.Machine, 0, len(machines))
		filter := d2hlp.NewStringFilter(names)
		for _, m := range machines {
			if filter.IsIncluded(m.Name) {
				res = append(res, m)
			}
		}
		machines = res

	}
	// get machine details
	return machines, nil
}

func (m *DRMAATracker) JobInfoFromMonitor(id string) (drmaa2interface.JobInfo, error) {
	if m.workloadManager != SonOfGridEngine {
		// TODO implement support for other systems
		return drmaa2interface.JobInfo{}, fmt.Errorf("unsupported workload manager")
	}
	jobStatus, err := gestatus.GetJob(id)
	if err != nil {
		// TODO: job might be finished or not existing?!?
		// TODO: check lookup table, which needs to be updated
		// in the background with qacct information
		return drmaa2interface.JobInfo{}, fmt.Errorf("failed to get job status with qstat -xml: %v", err)
	}
	ji := drmaa2interface.JobInfo{}
	// TODO need job state for everything
	state, err := QstatJobState(id) // yet another cli call...
	if err != nil && err.Error() == "does not exist" {
		// might be in failed state - only qacct knows...
		ji.State = drmaa2interface.Done
	} else {
		ji.State = ConvertQstatJobState(state)
	}
	// TODO job may be done or failed and not visible in qstat
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
