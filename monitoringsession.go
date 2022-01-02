package drmaa2os

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
)

// MonitoringSession implements a DRMAA2 monitoring session based on
// the JobTracker and Monitorer interface. Currently it is expected
// to have one object which implements both interfaces hence jobtracker
// and monitorer reference the same address.
type MonitoringSession struct {
	name       string
	jobtracker jobtracker.JobTracker
	monitorer  jobtracker.Monitorer
}

// CloseMonitoringSession disengages from the backend, i.e.
// closes potentially the connection and can lead to a non-usable
// monitoring session.
func (ms *MonitoringSession) CloseMonitoringSession() error {
	return ms.monitorer.CloseMonitoringSession(ms.name)
}

// GetAllJobs returns all visible jobs, consisting potentially
// multiple JobSession an external jobs. The returned jobs are
// potentially read-only and can't be manipulated (stopped, suspended).
// The filter can restrict the jobs returned, if no filter is required,
// drmaa2interface.CreateJobInfo() should be used as filter, which
// sets "Unset" values for all fields which are not nullable. See
// also GetAllJobsWithoutFilter().
func (ms *MonitoringSession) GetAllJobs(filter drmaa2interface.JobInfo) ([]drmaa2interface.Job, error) {
	ids, err := ms.monitorer.GetAllJobIDs(nil)
	if err != nil {
		return nil, fmt.Errorf("failed getting job list: %v", err)
	}
	jobs := make([]drmaa2interface.Job, 0, len(ids))
	for _, id := range ids {
		job := newMonitoringJob(id, ms.name, drmaa2interface.JobTemplate{}, ms.jobtracker, ms.monitorer)
		// TODO apply filter
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// GetAllQueues returns all queues. If filter is set to a list of strings,
// it only returns queue names which are defined in the filter.
func (ms *MonitoringSession) GetAllQueues(filter []string) ([]drmaa2interface.Queue, error) {
	queueNames, err := ms.monitorer.GetAllQueueNames(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue names: %v", err)
	}
	queues := make([]drmaa2interface.Queue, 0, len(queueNames))
	for _, queueName := range queueNames {
		queues = append(queues, drmaa2interface.Queue{
			Name: queueName,
		})
	}
	return queues, nil
}

// GetAllMachines returns all machines in the cluster. If the filter is
// set the result contains only existing machines defined by the filter.
func (ms *MonitoringSession) GetAllMachines(filter []string) ([]drmaa2interface.Machine, error) {
	return ms.monitorer.GetAllMachines(filter)
}

// GetAllReservations returns all advance(d) reservations. Currently not
// implemented.
func (ms *MonitoringSession) GetAllReservations() ([]drmaa2interface.Reservation, error) {
	return nil, fmt.Errorf("Reservations are currently unsupported")
}
