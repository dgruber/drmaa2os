package kubernetestracker

import (
	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func DRMAA2State(jc batchv1.JobInterface, jobid string) drmaa2interface.JobState {
	job, err := getJobByID(jc, jobid)
	if err != nil {
		return drmaa2interface.Undetermined
	}
	return convertJobStatus2JobState(&job.Status)
}
