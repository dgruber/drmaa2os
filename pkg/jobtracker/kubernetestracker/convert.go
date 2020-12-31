package kubernetestracker

import (
	"crypto/md5"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/api/batch/v1"
	k8sv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func volumeName(jobName, path string, kind string) string {
	sum := md5.Sum([]byte(path))
	return jobName + "-" + fmt.Sprintf("%.8x", sum) + "-" + kind + "-volume"
}

func pvcName(jobName, path string) string {
	sum := md5.Sum([]byte(path))
	return jobName + "-" + fmt.Sprintf("%.8x", sum) + "-pvc"
}

func configMapName(jobName, path string) string {
	sum := md5.Sum([]byte(path))
	return jobName + "-" + fmt.Sprintf("%.8x", sum) + "-configmap"
}

func secretName(jobName, path string) string {
	sum := md5.Sum([]byte(path))
	return jobName + "-" + fmt.Sprintf("%.8x", sum) + "-secret"
}

// newVolumes creates volumes for staging in data into the containers.
// The volumes are build on existing secrets or configmaps containing
// the data.
func newVolumes(jt drmaa2interface.JobTemplate) ([]k8sv1.Volume, error) {

	if jt.StageInFiles != nil {
		// naming scheme of the objects jobname-filename-{secret|cm}-volume
		volumes := make([]k8sv1.Volume, 0, 2)
		for path, v := range jt.StageInFiles {
			if strings.HasPrefix(v, "secret-data:") {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "secret"),
						VolumeSource: k8sv1.VolumeSource{
							Secret: &k8sv1.SecretVolumeSource{
								SecretName: secretName(jt.JobName, path),
							}}})
			} else if strings.HasPrefix(v, "configmap-data:") {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "cm"),
						VolumeSource: k8sv1.VolumeSource{
							ConfigMap: &k8sv1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: configMapName(jt.JobName, path)},
							}}})
			} else if strings.HasPrefix(v, "hostpath:") {
				sourcePath := strings.TrimPrefix(v, "hostpath:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "hostpath"),
						VolumeSource: k8sv1.VolumeSource{
							HostPath: &k8sv1.HostPathVolumeSource{
								Path: sourcePath,
							},
						}})

			} else if strings.HasPrefix(v, "configmap:") {
				existingConfigMapName := strings.TrimPrefix(v, "configmap:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "cm"),
						VolumeSource: k8sv1.VolumeSource{
							ConfigMap: &k8sv1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: existingConfigMapName},
							}}})
			} else if strings.HasPrefix(v, "secret:") {
				existingSecretName := strings.TrimPrefix(v, "secret:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "secret"),
						VolumeSource: k8sv1.VolumeSource{
							Secret: &k8sv1.SecretVolumeSource{
								SecretName: existingSecretName,
							}}})
			} else if strings.HasPrefix(v, "pvc:") {
				existingPVCName := strings.TrimPrefix(v, "pvc:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "pvc"),
						VolumeSource: k8sv1.VolumeSource{
							PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{
								ClaimName: existingPVCName,
							}}})
			} else if strings.HasPrefix(v, "gce-disk:") {
				existingPDName := strings.TrimPrefix(v, "gce-disk:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "gce-disk"),
						VolumeSource: k8sv1.VolumeSource{
							GCEPersistentDisk: &k8sv1.GCEPersistentDiskVolumeSource{
								PDName:   existingPDName,
								FSType:   "ext4",
								ReadOnly: false,
							}}})
			} else if strings.HasPrefix(v, "gce-disk-read:") {
				existingPDName := strings.TrimPrefix(v, "gce-disk-read:")
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "gce-disk-read"),
						VolumeSource: k8sv1.VolumeSource{
							GCEPersistentDisk: &k8sv1.GCEPersistentDiskVolumeSource{
								PDName:   existingPDName,
								FSType:   "ext4",
								ReadOnly: true,
							}}})
			} else if strings.HasPrefix(v, "storageclass:") {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "storageclass"),
						VolumeSource: k8sv1.VolumeSource{
							PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{
								ClaimName: pvcName(jt.JobName, path),
							},
						}})
			}
		}
		return volumes, nil
	}
	return nil, nil
}

func getVolumeMounts(jt drmaa2interface.JobTemplate) []v1.VolumeMount {
	if len(jt.StageInFiles) == 0 {
		return nil
	}
	vmounts := make([]v1.VolumeMount, 0, len(jt.StageInFiles))
	for k, v := range jt.StageInFiles {
		_, file := filepath.Split(k)
		if strings.HasPrefix(v, "secret-data:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "secret"),
				MountPath: k,
				SubPath:   file,
			})
		} else if strings.HasPrefix(v, "configmap-data:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "cm"),
				MountPath: k,
				SubPath:   file,
			})
		} else if strings.HasPrefix(v, "hostpath:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "hostpath"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, "configmap:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "cm"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, "secret:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "secret"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, "pvc:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pvc"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, "gce-disk:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pd"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, "gce-disk-read:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pd"),
				MountPath: k,
				ReadOnly:  true,
			})
		} else if strings.HasPrefix(v, "storageclass:") {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "storageclass"),
				MountPath: k,
			})
		}
	}
	return vmounts
}

func newContainers(jt drmaa2interface.JobTemplate) ([]k8sv1.Container, error) {
	if jt.JobCategory == "" {
		return nil, errors.New("JobCategory (container image name) not set in JobTemplate")
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

	c.VolumeMounts = getVolumeMounts(jt)

	for name, value := range jt.JobEnvironment {
		c.Env = append(c.Env, k8sv1.EnvVar{Name: name, Value: value})
	}

	// spec.template.spec.containers[0].name: Required value"
	if jt.JobName == "" {
		c.Name = "drmaa2os"
	}

	// if len(jt.CandidateMachines) == 1 {
	//	c = jt.CandidateMachines[0]
	// }
	return []k8sv1.Container{c}, nil
}

func newNodeSelector(jt drmaa2interface.JobTemplate) (map[string]string, error) {
	return nil, nil
}

/* 	deadlineTime returns the deadline of the job as int pointer converting from
    AbsoluteTime to a relative time.
	"
	Specifies a deadline after which the implementation or the DRM system SHOULD change the job state to
		any of the “Terminated” states (see Section 8.1).
    	The support for this attribute is optional, as expressed by the
       	- DrmaaCapability::JT_DEADLINE
		DeadlineTime is defined as AbsoluteTime.
	"
*/
func deadlineTime(jt drmaa2interface.JobTemplate) (*int64, error) {
	var deadline int64
	if !jt.DeadlineTime.IsZero() {
		if jt.DeadlineTime.After(time.Now()) {
			deadline = jt.DeadlineTime.Unix() - time.Now().Unix()
		} else {
			return nil, fmt.Errorf("deadlineTime (%s) in job template is in the past", jt.DeadlineTime.String())
		}
	}
	return &deadline, nil
}

// https://godoc.org/k8s.io/api/core/v1#PodSpec
// https://github.com/kubernetes/kubernetes/blob/886e04f1fffbb04faf8a9f9ee141143b2684ae68/pkg/api/types.go
func newPodSpec(v []k8sv1.Volume, c []k8sv1.Container, ns map[string]string, activeDeadline *int64) k8sv1.PodSpec {
	spec := k8sv1.PodSpec{
		Volumes:       v,
		Containers:    c,
		NodeSelector:  ns,
		RestartPolicy: k8sv1.RestartPolicyNever,
	}
	if *activeDeadline > 0 {
		spec.ActiveDeadlineSeconds = activeDeadline
	}
	return spec
}

func addExtensions(job *batchv1.Job, jt drmaa2interface.JobTemplate) *batchv1.Job {
	if jt.ExtensionList == nil {
		return job
	}
	if labels, set := jt.ExtensionList["labels"]; set && labels != "" {
		// "key=value,key=value,..."
		for _, label := range strings.Split(labels, ",") {
			l := strings.Split(label, "=")
			if len(l) == 2 {
				if l[0] == "drmaa2jobsession" {
					continue // don't allow to override job session
				}
				job.Labels[l[0]] = l[1]
			}
		}
	}

	if scheduler, set := jt.ExtensionList["scheduler"]; set && scheduler != "" {
		job.Spec.Template.Spec.SchedulerName = scheduler
	}

	if privileged, set := jt.ExtensionList["privileged"]; set && privileged != "" {
		if strings.ToUpper(privileged) == "TRUE" {
			for i := range job.Spec.Template.Spec.Containers {
				privileged := true
				job.Spec.Template.Spec.Containers[i].SecurityContext = &v1.SecurityContext{
					Privileged: &privileged,
				}
			}
		}
	}

	return job
}

func convertJob(jobsession, namespace string, jt drmaa2interface.JobTemplate) (*batchv1.Job, error) {
	volumes, err := newVolumes(jt)
	if err != nil {
		return nil, fmt.Errorf("converting job (newVolumes): %s", err)
	}
	containers, err := newContainers(jt)
	if err != nil {
		return nil, fmt.Errorf("converting job (newContainer): %s", err)
	}
	nodeSelector, err := newNodeSelector(jt)
	if err != nil {
		return nil, fmt.Errorf("converting job (newNodeSelector): %s", err)
	}

	// settings for command etc.
	dl, err := deadlineTime(jt)
	if err != nil {
		return nil, err
	}
	podSpec := newPodSpec(volumes, containers, nodeSelector, dl)

	var one int32 = 1
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		// Standard object's metadata.
		// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
		// +optional
		ObjectMeta: metav1.ObjectMeta{
			Name:         jt.JobName,
			Labels:       map[string]string{"drmaa2jobsession": jobsession},
			GenerateName: "drmaa2os",
			Namespace:    namespace,
		},
		// Specification of the desired behavior of a job.
		// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
		// +optional
		Spec: batchv1.JobSpec{
			/*ManualSelector: ,
			Selector: &unversioned.LabelSelector{
				MatchLabels: options.labels,
			}, */
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &one,

			// Describes the pod that will be created when executing a job.
			// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
			Template: k8sv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:         "drmaa2osjob",
					GenerateName: "drmaa2os",
					//Labels: options.labels,
					Namespace: namespace,
				},
				Spec: podSpec,
			},
		},
	}
	return addExtensions(&job, jt), nil
}
