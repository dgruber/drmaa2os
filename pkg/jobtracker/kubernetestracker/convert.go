package kubernetestracker

import (
	"errors"
	"fmt"
	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/api/batch/v1"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newVolumes(jt drmaa2interface.JobTemplate) ([]k8sv1.Volume, error) {
	//v := k8sv1.Volume{}
	return nil, nil
}

func newContainers(jt drmaa2interface.JobTemplate) ([]k8sv1.Container, error) {
	if jt.JobCategory == "" {
		return nil, errors.New("JobCategory (image name) not set in JobTemplate")
	}
	if jt.RemoteCommand == "" {
		return nil, errors.New("RemoteCommand not set in JobTemplate")
	}
	c := k8sv1.Container{
		Name:       jt.JobName,
		Image:      jt.JobCategory,
		Command:    []string{jt.RemoteCommand},
		Args:       jt.Args,
		WorkingDir: jt.WorkingDirectory,
	}
	// if len(jt.CandidateMachines) == 1 {
	//	c = jt.CandidateMachines[0]
	// }
	return []k8sv1.Container{c}, nil
}

func newNodeSelector(jt drmaa2interface.JobTemplate) (map[string]string, error) {
	return nil, nil
}

// https://github.com/kubernetes/kubernetes/blob/886e04f1fffbb04faf8a9f9ee141143b2684ae68/pkg/api/types.go
func newPodSpec(v []k8sv1.Volume, c []k8sv1.Container, ns map[string]string) k8sv1.PodSpec {
	return k8sv1.PodSpec{
		Volumes:      v,
		Containers:   c,
		NodeSelector: ns,
	}
}

func convertJob(jt drmaa2interface.JobTemplate) (*batchv1.Job, error) {

	volumes, err := newVolumes(jt)
	if err != nil {
		return nil, fmt.Errorf("error converting job (newVolumes): %s", err)
	}

	containers, err := newContainers(jt)
	if err != nil {
		return nil, fmt.Errorf("error converting job (newContainer): %s", err)
	}

	nodeSelector, err := newNodeSelector(jt)
	if err != nil {
		return nil, fmt.Errorf("error converting job (newNodeSelector): %s", err)
	}

	podSpec := newPodSpec(volumes, containers, nodeSelector)

	var one int32 = 1
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		// Standard object's metadata.
		// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
		// +optional
		ObjectMeta: metav1.ObjectMeta{
			Name: jt.JobName,
			//Namespace: v1.NamespaceDefault,
			//Labels: options.labels,
		},
		// Specification of the desired behavior of a job.
		// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
		// +optional
		Spec: batchv1.JobSpec{
			/*ManualSelector: ,
			Selector: &unversioned.LabelSelector{
				MatchLabels: options.labels,
			}, */
			Parallelism: &one,
			Completions: &one,

			// Describes the pod that will be created when executing a job.
			// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
			Template: k8sv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: jt.JobName,
					//Labels: options.labels,
				},
				Spec: podSpec,
			},
		},
	}, nil
}
