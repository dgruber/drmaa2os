package dockertracker

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

// Implements the Monitorer interface on top of the JobTracker interface
// so that MonitoringSessions can be created.

func (dt *DockerTracker) OpenMonitoringSession(name string) error {
	return nil
}

func (dt *DockerTracker) CloseMonitoringSession(name string) error {
	return nil
}

func (dt *DockerTracker) GetAllJobIDs(filter *drmaa2interface.JobInfo) ([]string, error) {
	if err := dt.check(); err != nil {
		return nil, err
	}
	// TODO convert jobinfo filter to container filter
	if filter != nil {
		return nil, fmt.Errorf("job filter not implemented")
	}
	f := filters.NewArgs()
	containers, err := dt.cli.ContainerList(context.Background(),
		container.ListOptions{Filters: f, All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %v", err)
	}
	ids := make([]string, 0, len(containers))
	for i := range containers {
		ids = append(ids, containers[i].ID)
	}
	return ids, nil
}

func (dt *DockerTracker) GetAllQueueNames(filter []string) ([]string, error) {
	return []string{}, nil
}

func (dt *DockerTracker) GetAllMachines(filter []string) ([]drmaa2interface.Machine, error) {
	info, err := dt.cli.Info(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get docker host info: %v", err)
	}

	for f := range filter {
		if filter[f] == info.Name {
			return []drmaa2interface.Machine{}, nil
		}
	}

	var arch drmaa2interface.CPU
	if info.Architecture == "x86_64" {
		arch = drmaa2interface.IA64
	}

	var os drmaa2interface.OS
	if info.OSType == "linux" {
		os = drmaa2interface.Linux
	}

	return []drmaa2interface.Machine{
		{
			Name:           info.Name,
			Architecture:   arch,
			Sockets:        1, // do we know better?
			CoresPerSocket: int64(info.NCPU),
			ThreadsPerCore: 1, // do we know better?
			PhysicalMemory: info.MemTotal,
			VirtualMemory:  info.MemTotal,
			OS:             os,
		},
	}, nil
}

// JobInfoFromMonitor might collect job state and job info in a
// different way as a JobSession with persistent storage does
func (dt *DockerTracker) JobInfoFromMonitor(id string) (ji drmaa2interface.JobInfo, err error) {
	if err := dt.check(); err != nil {
		return ji, err
	}
	container, err := dt.cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return ji, err
	}
	return containerToDRMAA2JobInfo(container)
}
