package kubernetestracker

import (
	"errors"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	k8sapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientBatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

func jobStateChange(jc clientBatchv1.JobInterface, job *batchv1.Job, action string) error {
	if jc == nil || job == nil {
		return errors.New("can't change job status: job is nil")
	}
	switch action {
	case "suspend":
		return errors.New("Unsupported Operation")
	case "resume":
		return errors.New("Unsupported Operation")
	case "hold":
		return errors.New("Unsupported Operation")
	case "release":
		return errors.New("Unsupported Operation")
	case "terminate":
		return jc.Delete(job.GetName(), &k8sapi.DeleteOptions{})
	}
	return fmt.Errorf("Undefined job operation")
}
