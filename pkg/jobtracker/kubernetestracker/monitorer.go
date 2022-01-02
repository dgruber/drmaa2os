package kubernetestracker

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgruber/drmaa2interface"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesTracker implements the Monitorer interface on top of the
// JobTracker interface so that MonitoringSessions can be created by
// the SessionManager.

func (kt *KubernetesTracker) OpenMonitoringSession(name string) error {
	return nil
}

func (kt *KubernetesTracker) CloseMonitoringSession(name string) error {
	return nil
}

func (kt *KubernetesTracker) GetAllJobIDs(filter *drmaa2interface.JobInfo) ([]string, error) {
	jc, err := getJobsClient(kt.clientSet, kt.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get all job IDs: %s", err.Error())
	}
	// TODO implement filter
	jobsList, err := jc.List(context.TODO(), k8sapi.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing jobs with client: %s", err.Error())
	}
	ids := make([]string, 0, len(jobsList.Items))
	for _, job := range jobsList.Items {
		ids = append(ids, string(job.Name))
	}
	return ids, nil
}

// GetAllQueueNames returns all namespaces. If filter is != nil it returns
// only existing namespaces which are defined by the filter.
func (kt *KubernetesTracker) GetAllQueueNames(filter []string) ([]string, error) {
	// queues are namespaces
	namespaceList, err := kt.clientSet.CoreV1().Namespaces().List(context.Background(),
		k8sapi.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed get kubernetes namespace list: %v", err)
	}
	namespaces := make([]string, 0, len(namespaceList.Items))
	for _, namespace := range namespaceList.Items {
		if filter != nil {
			for f := range filter {
				if namespace.Name == filter[f] {
					namespaces = append(namespaces, namespace.Name)
				}
			}
		} else {
			// unfiltered
			namespaces = append(namespaces, namespace.Name)
		}
	}
	return namespaces, nil
}

func (kt *KubernetesTracker) GetAllMachines(filter []string) ([]drmaa2interface.Machine, error) {
	nodeList, err := kt.clientSet.CoreV1().Nodes().List(context.Background(),
		k8sapi.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed get kubernetes node list: %v", err)
	}
	machines := make([]drmaa2interface.Machine, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {

		filtered := false
		for _, f := range filter {
			if node.Name == f {
				filtered = true
				break
			}
		}
		if filtered {
			continue
		}

		mem, _ := node.Status.Capacity.Memory().AsInt64()
		var os drmaa2interface.OS
		if node.Status.NodeInfo.OperatingSystem == "linux" {
			os = drmaa2interface.Linux
		} else if node.Status.NodeInfo.OperatingSystem == "darwin" {
			os = drmaa2interface.MacOS
		} else if node.Status.NodeInfo.OperatingSystem == "windows" {
			os = drmaa2interface.Win
		}

		var arch drmaa2interface.CPU
		if strings.HasPrefix(node.Status.NodeInfo.Architecture, "ppc64") {
			arch = drmaa2interface.PowerPC64
		} else if node.Status.NodeInfo.Architecture == "amd64" {
			arch = drmaa2interface.IA64
		} else if strings.HasPrefix(node.Status.NodeInfo.Architecture, "arm64") {
			arch = drmaa2interface.ARM64
		}

		osver := strings.Split(node.Status.NodeInfo.KernelVersion, ".")
		if len(osver) < 3 {
			osver = append(osver, "0")
			osver = append(osver, "0")
			osver = append(osver, "0")
		}
		cores, _ := node.Status.Capacity.Cpu().AsInt64()
		machines = append(machines, drmaa2interface.Machine{
			Name:           node.Name,
			PhysicalMemory: mem,
			VirtualMemory:  mem,
			OS:             os,
			Architecture:   arch,
			OSVersion: drmaa2interface.Version{
				Major: osver[0],
				Minor: osver[1] + "." + osver[2],
			},
			Sockets:        1, // don't know better
			CoresPerSocket: cores,
			ThreadsPerCore: 1, // don't know better
		})
	}
	return machines, nil
}

// JobInfoFromMonitor might collect job state and job info in a
// different way as a JobSession with persistent storage does
func (kt *KubernetesTracker) JobInfoFromMonitor(id string) (ji drmaa2interface.JobInfo, err error) {
	return kt.JobInfo(id)
}
