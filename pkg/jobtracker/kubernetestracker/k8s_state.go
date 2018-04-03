package kubernetestracker

import (
	"github.com/dgruber/drmaa2interface"
	"k8s.io/api/batch/v1"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func convertJobStatus2JobState(status *v1.JobStatus) drmaa2interface.JobState {
	if status == nil {
		return drmaa2interface.Undetermined
	}
	// https://kubernetes.io/docs/api-reference/batch/v1/definitions/#_v1_jobstatus
	if status.Succeeded >= 1 {
		return drmaa2interface.Done
	}
	if status.Failed >= 1 {
		return drmaa2interface.Failed
	}
	if status.Active >= 1 {
		return drmaa2interface.Running
	}
	return drmaa2interface.Undetermined
}

func DRMAA2State(jc batchv1.JobInterface, jobid string) drmaa2interface.JobState {
	job, err := getJobByID(jc, jobid)
	if err != nil {
		return drmaa2interface.Undetermined
	}
	return convertJobStatus2JobState(&job.Status)
}

func exitStatusFromJobState(status drmaa2interface.JobState) int {
	switch status {
	case drmaa2interface.Failed:
		return 1
	case drmaa2interface.Done:
		return 0
	}
	return 0
}

func JobToJobInfo(jc batchv1.JobInterface, jobid string) (drmaa2interface.JobInfo, error) {
	job, err := getJobByID(jc, jobid)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}

	ji := drmaa2interface.JobInfo{}

	ji.SubmissionTime = job.CreationTimestamp.Time
	if job.Status.StartTime != nil {
		ji.DispatchTime = job.Status.StartTime.Time
	}
	if job.Status.CompletionTime != nil {
		ji.FinishTime = job.Status.CompletionTime.Time
	}
	ji.State = convertJobStatus2JobState(&job.Status)
	ji.ID = jobid
	ji.ExitStatus = exitStatusFromJobState(ji.State)

	return ji, nil
}
