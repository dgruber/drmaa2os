package kubernetestracker

import (
	"github.com/dgruber/drmaa2interface"
	"k8s.io/client-go/kubernetes"
)

func DRMAA2State(cs *kubernetes.Clientset, jobid string) drmaa2interface.JobState {
	job, err := getJobByID(cs, jobid)
	if err != nil {
		return drmaa2interface.Undetermined
	}
	return convertJobStatus2JobState(&job.Status)
}
