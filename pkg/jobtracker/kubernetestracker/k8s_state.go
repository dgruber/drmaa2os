package kubernetestracker

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgruber/drmaa2interface"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
	if status.CompletionTime != nil && status.CompletionTime.Time.Before(time.Now()) {
		return drmaa2interface.Failed
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

// JobToJobInfo converts a kubernetes job to a DRMAA2 JobInfo representation.
func JobToJobInfo(jc batchv1.JobInterface, jobid string) (drmaa2interface.JobInfo, error) {
	ji := drmaa2interface.JobInfo{}
	job, err := getJobByID(jc, jobid)
	if err != nil {
		return ji, err
	}
	ji.Slots = 1
	ji.SubmissionTime = job.CreationTimestamp.Time
	if job.Status.StartTime != nil {
		ji.DispatchTime = job.Status.StartTime.Time
	}
	if job.Status.CompletionTime != nil {
		ji.FinishTime = job.Status.CompletionTime.Time
		ji.WallclockTime = ji.FinishTime.Sub(ji.DispatchTime)
	}
	ji.State = convertJobStatus2JobState(&job.Status)
	ji.ID = jobid
	ji.ExitStatus = exitStatusFromJobState(ji.State)
	return ji, nil
}

// GetJobOutput returns the output of a job pod after after it has been finished.
func GetJobOutput(cs kubernetes.Interface, namespace string, jobID string) ([]byte, error) {
	podList, err := cs.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "job-name=" + jobID,
	})
	if err != nil || len(podList.Items) <= 0 {
		return nil, fmt.Errorf("could not get pods of job %s in namespace %s: %v",
			jobID, namespace, err)
	}

	podName := podList.Items[0].Name

	if len(podList.Items) != 1 {
		// might be a problem with the container pod and there are restarts
		// find the last pod which has been created
		last := time.Unix(0, 0)
		var lastPod *corev1.Pod
		for _, pod := range podList.Items {
			if pod.CreationTimestamp.Time.After(last) {
				last = pod.CreationTimestamp.Time
				lastPod = &pod
			}
		}
		podName = lastPod.Name
	}

	req := cs.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: jobID,
		Follow:    false,
	})
	output, err := req.Stream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get output stream of pod %s from job %s in namespace %s: %v",
			podList.Items[0].Name, jobID, namespace, err)
	}
	return ioutil.ReadAll(output)
}
