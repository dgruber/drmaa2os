package kubernetestracker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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
	for _, condition := range status.Conditions {
		if condition.Type == v1.JobFailed && condition.Status == corev1.ConditionTrue {
			return drmaa2interface.Failed
		}
		if condition.Type == v1.JobComplete && condition.Status == corev1.ConditionTrue {
			return drmaa2interface.Done
		}
	}
	// TODO: check for suspended state
	// From Kubernetes code base:
	// "The latest available observations of an object's current state. When a Job
	// fails, one of the conditions will have type "Failed" and status true. When
	// a Job is suspended, one of the conditions will have type "Suspended" and
	// status true; when the Job is resumed, the status of this condition will
	// become false. When a Job is completed, one of the conditions will have
	// type "Complete" and status true."
	if status.Succeeded >= 1 {
		return drmaa2interface.Done
	}

	if status.Failed >= 1 {
		return drmaa2interface.Failed
	}

	if status.Active >= 1 {
		return drmaa2interface.Running
	}
	// completed already
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
	ji.AllocatedMachines = []string{job.Spec.Template.Spec.NodeName}

	if IsDeadlineTimeException(job.Status.Conditions) {
		ji.SubState = "DeadlineExceeded"
	}

	return ji, nil
}

func IsComplete(c *[]v1.JobCondition) bool {
	for _, condition := range *c {
		if condition.Type == v1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsDeadlineTimeException(c []v1.JobCondition) bool {
	for _, condition := range c {
		if strings.Contains(condition.Reason, "DeadlineExceeded") {
			return true
		}
	}
	return false
}

// GetJobOutput returns the output of a job pod after after it has been finished.
func GetJobOutput(cs kubernetes.Interface, namespace string, jobID, podName string) ([]byte, error) {
	req := cs.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: jobID,
		Follow:    false,
	})
	output, err := req.Stream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get output stream of pod %s from job %s in namespace %s: %v",
			podName, jobID, namespace, err)
	}
	return ioutil.ReadAll(output)
}

func GetMachineNameForPod(cs kubernetes.Interface, namespace, podName string) (string, error) {
	pod, err := cs.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get pod %s in namespace %s: %v", podName, namespace, err)
	}

	return pod.Spec.NodeName, nil
}

// GetExitStatusOfJobContainer returns the exit status of a job container.
func GetExitStatusOfJobContainer(cs kubernetes.Interface, namespace, podName string) (int32, int32, string, error) {
	pod, err := cs.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to get pod %s in namespace %s: %v", podName, namespace, err)
	}
	message := pod.Status.Reason
	// assume: container 0 is the job container (this depends on the pod spec in the job spec)
	if len(pod.Status.ContainerStatuses) >= 1 && pod.Status.ContainerStatuses[0].State.Terminated != nil {
		exitCode := pod.Status.ContainerStatuses[0].State.Terminated.ExitCode
		terminationSignal := pod.Status.ContainerStatuses[0].State.Terminated.Signal
		//message := pod.Status.ContainerStatuses[0].State.Terminated.Message
		return exitCode, terminationSignal, message, nil
	}
	return 0, 0, message, fmt.Errorf("no container terminated status found for pod %s in namespace %s", podName, namespace)
}

func GetPodsForJob(cs kubernetes.Interface, namespace, jobID string) ([]corev1.Pod, error) {
	podList, err := cs.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "job-name=" + jobID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get pods of job %s in namespace %s: %v",
			jobID, namespace, err)
	}
	if len(podList.Items) <= 0 {
		return nil, fmt.Errorf("no active pod for job with label selector job-name=%s", jobID)
	}
	return podList.Items, nil
}

func GetLastStartedPod(pods []corev1.Pod) corev1.Pod {
	last := time.Unix(0, 0)
	var lastPod corev1.Pod
	for _, pod := range pods {
		log.Printf("found pod %v\n", pod.Name)
		if pod.CreationTimestamp.Time.After(last) {
			last = pod.CreationTimestamp.Time
			lastPod = pod
		}
	}
	return lastPod
}

// GetFirstPod is required when a job is deleted by deadline time.
// Here we have no control about restarts (seeing 3) as deadling takes
// precedence over backoff limits in Kubernetes. So, when deadline time
// is used there is no guarantee that we only have one pod.
func GetFirstPod(pods []corev1.Pod) corev1.Pod {
	return pods[0]
}
