package kubernetestracker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/extension"
	"github.com/dgruber/drmaa2os/pkg/helper"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const K8S_JT_EXTENSION_NAMESPACE = "namespace"
const K8S_JT_EXTENSION_LABELS = "labels"

type KubernetesTracker struct {
	clientSet  *kubernetes.Clientset
	jobsession string
	namespace  string
}

// init registers the Kubernetes job tracker at the SessionManager
func init() {
	drmaa2os.RegisterJobTracker(drmaa2os.KubernetesSession, NewAllocator())
}

type allocator struct{}

func NewAllocator() *allocator {
	return &allocator{}
}

// KubernetesTrackerParameters can be used as parameter in
// NewKubernetesSessionManager. Note, that the namespace
// if set must exist. If not set the "default" namespace
// is used.
type KubernetesTrackerParameters struct {
	Namespace string // if not set it will become "default"
	ClientSet *kubernetes.Clientset
}

// New is called by the SessionManager when a new JobSession is allocated.
// jobTrackerInitParams must be either a *kubernetes.Clientset or a
// KubernetesTrackerParameters struct or nil. If nil or KubernetesTrackerParameters
// has a nil clientset a new Kubernetes clientset is allocated.
func (a *allocator) New(jobSessionName string, jobTrackerInitParams interface{}) (jobtracker.JobTracker, error) {
	if jobTrackerInitParams != nil {
		switch v := jobTrackerInitParams.(type) {
		case *kubernetes.Clientset:
			return New(jobSessionName, "default", v)
		case KubernetesTrackerParameters:
			return New(jobSessionName, v.Namespace, v.ClientSet)
		default:
			return nil, errors.New("jobTrackerInitParams is not of type *kubernetes.Clientset or KubernetesTrackerParameters")

		}
	}
	return New(jobSessionName, "default", nil)
}

// New creates a new KubernetesTracker either by using a given kubernetes Clientset
// or by allocating a new one (if the parameter is zero).
func New(jobsession string, namespace string, cs *kubernetes.Clientset) (*KubernetesTracker, error) {
	if cs == nil {
		var err error
		cs, err = NewClientSet()
		if err != nil {
			return nil, err
		}
	}
	if namespace == "" {
		namespace = "default"
	}
	return &KubernetesTracker{
		clientSet:  cs,
		jobsession: jobsession,
		namespace:  namespace,
	}, nil
}

// ListJobCategories returns all container images which are currently
// found in the cluster. That does not mean that other container images
// can not be used.
func (kt *KubernetesTracker) ListJobCategories() ([]string, error) {
	nodeList, err := kt.clientSet.CoreV1().Nodes().List(context.Background(),
		k8sapi.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed get kubernetes node list: %v", err)
	}
	images := make([]string, 0, len(nodeList.Items))
	for i := range nodeList.Items {
		for j := range nodeList.Items[i].Status.Images {
			images = append(images, nodeList.Items[i].Status.Images[j].Names...)
		}
	}
	return images, nil
}

// ListJobs returns a list of job IDs associated with the current
// DRMAA2 job session.
func (kt *KubernetesTracker) ListJobs() ([]string, error) {
	jc, err := getJobsClient(kt.clientSet, kt.namespace)
	if err != nil {
		return nil, fmt.Errorf("ListJobs: %s", err.Error())
	}
	labelSelector := fmt.Sprintf("drmaa2jobsession=%s", kt.jobsession)
	jobsList, err := jc.List(context.TODO(), k8sapi.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("listing jobs with client: %s", err.Error())
	}
	ids := make([]string, 0, len(jobsList.Items))
	for _, job := range jobsList.Items {
		ids = append(ids, string(job.Name))
	}
	return ids, nil
}

// AddJob converts the given DRMAA2 job template into a batchv1.Job and creates
// the job within Kubernetes.
func (kt *KubernetesTracker) AddJob(jt drmaa2interface.JobTemplate) (string, error) {
	// unique job name is required for secrets, configmap names and pod
	if jt.JobName == "" {
		jt.JobName = fmt.Sprintf("d2-%d", time.Now().UnixNano())
	}

	// create secrets and configmaps with file contents from StageInFiles
	secrets, err := getJobStageInSecrets(jt)
	if err != nil {
		return "", err
	}
	for _, secret := range secrets {
		_, err := kt.clientSet.CoreV1().Secrets(kt.namespace).Create(context.TODO(),
			secret, k8sapi.CreateOptions{})
		if err != nil {
			return "", err
		}
	}

	configmaps, err := getJobStageInConfigMaps(jt)
	if err != nil {
		return "", err
	}
	for _, configmap := range configmaps {
		_, err := kt.clientSet.CoreV1().ConfigMaps(kt.namespace).Create(context.TODO(),
			configmap, k8sapi.CreateOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create configmap: %v", err)
		}
	}

	// create PVC when storage class is used
	pvcs, err := getJobStageInPVCs(jt)
	if err != nil {
		return "", err
	}
	for _, pvc := range pvcs {
		_, err := kt.clientSet.CoreV1().PersistentVolumeClaims(kt.namespace).Create(context.Background(),
			pvc, k8sapi.CreateOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create PVC for StorageClass: %w", err)
		}
	}

	job, err := convertJob(kt.jobsession, kt.namespace, jt)
	if err != nil {
		removeArtifacts(kt.clientSet, jt, kt.namespace)
		return "", fmt.Errorf("converting job template into a k8s job: %s", err.Error())
	}

	jc, err := getJobsClient(kt.clientSet, kt.namespace)
	if err != nil {
		removeArtifacts(kt.clientSet, jt, kt.namespace)
		return "", fmt.Errorf("get client: %s", err.Error())
	}
	j, err := jc.Create(context.Background(), job, k8sapi.CreateOptions{})
	if err != nil {
		removeArtifacts(kt.clientSet, jt, kt.namespace)
		return "", fmt.Errorf("failed creating new job: %s", err.Error())
	}

	// store JobTemplate in ConfigMap
	err = storeJobTemplateInConfigMap(kt.clientSet, jt, kt.namespace)

	return string(j.Name), err
}

func (kt *KubernetesTracker) AddArrayJob(jt drmaa2interface.JobTemplate, begin int, end int, step int, maxParallel int) (string, error) {
	return helper.AddArrayJobAsSingleJobs(jt, kt, begin, end, step)
}

func (kt *KubernetesTracker) ListArrayJobs(id string) ([]string, error) {
	return helper.ArrayJobID2GUIDs(id)
}

func (kt *KubernetesTracker) JobState(jobID string) (drmaa2interface.JobState, string, error) {
	jc, err := getJobsClient(kt.clientSet, kt.namespace)
	if err != nil {
		return drmaa2interface.Undetermined, "", nil
	}
	return DRMAA2State(jc, jobID), "", nil
}

func (kt *KubernetesTracker) JobInfo(jobID string) (drmaa2interface.JobInfo, error) {
	jc, err := getJobsClient(kt.clientSet, kt.namespace)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}
	// JobInfo should return data staged out directly
	ji, err := JobToJobInfo(jc, jobID)
	if err != nil {
		return drmaa2interface.JobInfo{}, err
	}
	if ji.State == drmaa2interface.Done || ji.State == drmaa2interface.Failed {
		podList, errGetPods := GetPodsForJob(kt.clientSet, kt.namespace, jobID)
		if errGetPods != nil {
			// might be normal if pod already finished
			return ji, nil
		}
		//podName := GetLastStartedPod(podList).Name
		podName := GetFirstPod(podList).Name

		// read job output through logs
		output, err := GetJobOutput(kt.clientSet, kt.namespace, jobID, podName)
		if err == nil {
			if ji.ExtensionList == nil {
				ji.ExtensionList = make(map[string]string)
			}
			ji.ExtensionList[extension.JobInfoK8sJSessionJobOutput] = string(output)
		} else {
			fmt.Printf("error reading job output: %v\n", err)
		}

		machine, err := GetMachineNameForPod(kt.clientSet, kt.namespace, podName)
		if err != nil {
			ji.AllocatedMachines = []string{}
		} else {
			ji.AllocatedMachines = []string{machine}
		}

		// When a job has a deadline, then the job / pod will mark as failed
		// but does not have the termination status "terminated" struct set.
		// Hence we need to retry here...
		exitCode, terminationSignal, message, err := GetExitStatusOfJobContainer(kt.clientSet, kt.namespace, podName)
		if err != nil {
			// deadline - no terminated struct
			if message == "DeadlineExceeded" {
				ji.ExitStatus = int(138)
				ji.TerminatingSignal = fmt.Sprintf("%d", 9)
				ji.SubState = message
			}
		} else {
			ji.ExitStatus = int(exitCode)
			ji.TerminatingSignal = fmt.Sprintf("%d", terminationSignal)
			ji.SubState = message
		}
		// when job is finished the sidecar of the job
		// triggers an "epilog" job / or stores data in a
		// config map - this data needs to be read and
		// put into the JobInfo object
		/*
			bytesOutput, err := getConfigMapData(kt.clientSet, kt.namespace, jobID+"-output-configmap", "output")
			if err != nil {
				return ji, nil
			}
			if ji.ExtensionList == nil {
				ji.ExtensionList = make(map[string]string)
			}
			ji.ExtensionList["sidecar-output"] = string(bytesOutput)
		*/
	}
	return ji, nil
}

// JobControl changes the state of the given job by execution the given action
// (suspend, resume, hold, release, terminate).
func (kt *KubernetesTracker) JobControl(jobid, state string) error {
	jc, job, err := getJobInterfaceAndJob(kt.clientSet, jobid, kt.namespace)
	if err != nil {
		return fmt.Errorf("JobControl failed for jobID %s and action %s: %w",
			jobid, state, err)
	}
	return jobStateChange(jc, job, state)
}

// Wait returns when the job is in one of the given states or when a timeout
// occurs (errors then).
func (kt *KubernetesTracker) Wait(jobid string, timeout time.Duration, states ...drmaa2interface.JobState) error {
	return helper.WaitForState(kt, jobid, timeout, states...)
}

// DeleteJob removes a finished job and the objects created along
// with the job (like configmaps and secrets) Kubernetes.
func (kt *KubernetesTracker) DeleteJob(jobid string) error {
	jc, job, err := getJobInterfaceAndJob(kt.clientSet, jobid, kt.namespace)
	if err != nil {
		return fmt.Errorf("DeleteJob error: %w", err)
	}
	err = deleteJob(jc, job)
	if err != nil {
		return err
	}
	return removeArtifactsByJobID(kt.clientSet, jobid, kt.namespace)
}
