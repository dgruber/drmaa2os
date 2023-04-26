package kubernetestracker

import (
	"crypto/md5"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/extension"
	batchv1 "k8s.io/api/batch/v1"
	k8sv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
			if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsSecretB64Prefix) {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "secret"),
						VolumeSource: k8sv1.VolumeSource{
							Secret: &k8sv1.SecretVolumeSource{
								SecretName: secretName(jt.JobName, path),
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsConfigMapB64Prefix) {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "cm"),
						VolumeSource: k8sv1.VolumeSource{
							ConfigMap: &k8sv1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: configMapName(jt.JobName, path)},
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromHostPathPrefix) {
				sourcePath := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromHostPathPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "hostpath"),
						VolumeSource: k8sv1.VolumeSource{
							HostPath: &k8sv1.HostPathVolumeSource{
								Path: sourcePath,
							},
						}})

			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromConfigMapPrefix) {
				existingConfigMapName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromConfigMapPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "cm"),
						VolumeSource: k8sv1.VolumeSource{
							ConfigMap: &k8sv1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: existingConfigMapName},
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromSecretPrefix) {
				existingSecretName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromSecretPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "secret"),
						VolumeSource: k8sv1.VolumeSource{
							Secret: &k8sv1.SecretVolumeSource{
								SecretName: existingSecretName,
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromPVCPrefix) {
				existingPVCName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromPVCPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "pvc"),
						VolumeSource: k8sv1.VolumeSource{
							PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{
								ClaimName: existingPVCName,
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskPrefix) {
				existingPDName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "gce-disk"),
						VolumeSource: k8sv1.VolumeSource{
							GCEPersistentDisk: &k8sv1.GCEPersistentDiskVolumeSource{
								PDName:   existingPDName,
								FSType:   "ext4",
								ReadOnly: false,
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix) {
				existingPDName := strings.TrimPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix)
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "gce-disk-read"),
						VolumeSource: k8sv1.VolumeSource{
							GCEPersistentDisk: &k8sv1.GCEPersistentDiskVolumeSource{
								PDName:   existingPDName,
								FSType:   "ext4",
								ReadOnly: true,
							}}})
			} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromStorageClassNamePrefix) {
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "storageclass"),
						VolumeSource: k8sv1.VolumeSource{
							PersistentVolumeClaim: &k8sv1.PersistentVolumeClaimVolumeSource{
								ClaimName: pvcName(jt.JobName, path),
							},
						}})
			} else if strings.HasPrefix(v, "nfs:") {
				nfs := strings.Split(v, ":")
				if len(nfs) != 3 {
					return nil, errors.New("nfs source config needs to be in format nfs:server:path")
				}
				volumes = append(volumes,
					k8sv1.Volume{
						Name: volumeName(jt.JobName, path, "nfs"),
						VolumeSource: k8sv1.VolumeSource{
							NFS: &k8sv1.NFSVolumeSource{
								Server: nfs[1],
								Path:   nfs[2],
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
		if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsSecretB64Prefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "secret"),
				MountPath: k,
				SubPath:   file,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInAsConfigMapB64Prefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "cm"),
				MountPath: k,
				SubPath:   file,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromHostPathPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "hostpath"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromConfigMapPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "cm"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromSecretPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "secret"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromPVCPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pvc"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pd"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromGCEDiskReadOnlyPrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "pd"),
				MountPath: k,
				ReadOnly:  true,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromStorageClassNamePrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "storageclass"),
				MountPath: k,
			})
		} else if strings.HasPrefix(v, extension.JobTemplateK8sStageInFromNFSVolumePrefix) {
			vmounts = append(vmounts, v1.VolumeMount{
				Name:      volumeName(jt.JobName, k, "nfs"),
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

/*
	 	deadlineTime returns the deadline of the job as int pointer converting from
	    AbsoluteTime to a relative time.
		"
		Specifies a deadline after which the implementation or the DRM system SHOULD change the job state to
			any of the “Terminated” states (see Section 8.1).
	    	The support for this attribute is optional, as expressed by the
	       	- DrmaaCapability::JT_DEADLINE
			DeadlineTime is defined as AbsoluteTime.
		"
*/
func deadlineTime(jt drmaa2interface.JobTemplate) (int64, error) {
	var deadline int64
	deadline = -1 // unset

	if !jt.DeadlineTime.IsZero() {
		if jt.DeadlineTime.After(time.Now()) {
			deadline = jt.DeadlineTime.Unix() - time.Now().Unix()
		} else {
			return 0, fmt.Errorf("deadlineTime (%s) in job template is in the past",
				jt.DeadlineTime.String())
		}
	}
	return deadline, nil
}

// https://godoc.org/k8s.io/api/core/v1#PodSpec
// https://github.com/kubernetes/kubernetes/blob/886e04f1fffbb04faf8a9f9ee141143b2684ae68/pkg/api/types.go
func newPodSpec(v []k8sv1.Volume, c []k8sv1.Container, ns map[string]string) k8sv1.PodSpec {
	spec := k8sv1.PodSpec{
		Volumes:       v,
		Containers:    c,
		NodeSelector:  ns,
		RestartPolicy: k8sv1.RestartPolicyNever,
	}
	return spec
}

func addExtensions(job *batchv1.Job, jt drmaa2interface.JobTemplate) *batchv1.Job {
	if jt.ExtensionList == nil {
		return job
	}
	if labels, set := jt.ExtensionList[extension.JobTemplateK8sLabels]; set && labels != "" {
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

	if scheduler, set := jt.ExtensionList[extension.JobTemplateK8sScheduler]; set && scheduler != "" {
		job.Spec.Template.Spec.SchedulerName = scheduler
	}

	if privileged, set := jt.ExtensionList[extension.JobTemplateK8sPrivileged]; set && privileged != "" {
		if strings.ToUpper(privileged) == "TRUE" {
			for i := range job.Spec.Template.Spec.Containers {
				p := true
				fmt.Printf("add extension privileged=%s\n", privileged)
				if job.Spec.Template.Spec.Containers[i].SecurityContext == nil {
					job.Spec.Template.Spec.Containers[i].SecurityContext = &v1.SecurityContext{
						Privileged: &p,
					}
				} else {
					job.Spec.Template.Spec.Containers[i].SecurityContext.Privileged = &p
				}
			}
		}
	}

	if runasuser, set := jt.ExtensionList["runasuser"]; set && runasuser != "" {
		p, err := strconv.ParseInt(runasuser, 10, 64)
		if err != nil {
			fmt.Printf("runasuser: %s\n", err)
		} else {
			if job.Spec.Template.Spec.SecurityContext == nil {
				job.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
					RunAsUser: &p,
				}
			} else {
				job.Spec.Template.Spec.SecurityContext.RunAsUser = &p
			}
		}
	}

	if runasgroup, set := jt.ExtensionList["runasgroup"]; set && runasgroup != "" {
		p, err := strconv.ParseInt(runasgroup, 10, 64)
		if err != nil {
			fmt.Printf("runasgroup: %s\n", err)
		} else {
			if job.Spec.Template.Spec.SecurityContext == nil {
				job.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
					RunAsGroup: &p,
				}
			} else {
				job.Spec.Template.Spec.SecurityContext.RunAsGroup = &p
			}
		}
	}

	if fsuser, set := jt.ExtensionList["fsgroup"]; set && fsuser != "" {
		p, err := strconv.ParseInt(fsuser, 10, 64)
		if err != nil {
			fmt.Printf("fsgroup: %s\n", err)
		} else {
			if job.Spec.Template.Spec.SecurityContext == nil {
				job.Spec.Template.Spec.SecurityContext = &v1.PodSecurityContext{
					FSGroup: &p,
				}
			} else {
				job.Spec.Template.Spec.SecurityContext.FSGroup = &p
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

	podSpec := newPodSpec(volumes, containers, nodeSelector)

	// add enviornment variables from pre-existing secrets and config maps
	envFrom := []v1.EnvFromSource{}

	for k, v := range jt.ExtensionList {
		// both should work "env-from-secret" and "env-from-secrets"
		if strings.HasPrefix(k, extension.JobTemplateK8sEnvFromSecret) {
			for _, secret := range strings.Split(v, ":") {
				if secret == "" {
					continue
				}
				envFrom = append(envFrom,
					v1.EnvFromSource{
						SecretRef: &v1.SecretEnvSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: secret,
							},
						},
					})
			}
		}
		if strings.HasPrefix(k, extension.JobTemplateK8sEnvFromConfigMap) {
			for _, configmap := range strings.Split(v, ":") {
				if configmap == "" {
					continue
				}
				envFrom = append(envFrom,
					v1.EnvFromSource{
						ConfigMapRef: &v1.ConfigMapEnvSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: configmap,
							},
						},
					})
			}
		}
	}

	if envFrom != nil {
		podSpec.Containers[0].EnvFrom = envFrom
	}

	if jt.ExtensionList != nil && jt.ExtensionList["pullpolicy"] != "" {
		switch strings.ToLower(jt.ExtensionList["pullpolicy"]) {
		case "always":
			podSpec.Containers[0].ImagePullPolicy = v1.PullAlways
		case "never":
			podSpec.Containers[0].ImagePullPolicy = v1.PullNever
		case "ifnotpresent":
			podSpec.Containers[0].ImagePullPolicy = v1.PullIfNotPresent
		}
		// unknown pull policy will be ignored
	}

	// Add sidecar which stores the output of the job in a configmap.
	// This is not needed as the job output is read from the logs of
	// the pod object. But it is useful for storing the output in a
	// configmap which can be consumed by another job.
	// DRMAA2_JOB_OUTPUT_IN_JOBINFO is deprecated and will be renamed
	// in future versions. Please use the constant for drmaa2-basic-sidecar.
	if jt.ExtensionList != nil {
		jo, exists := jt.ExtensionList[extension.JobTemplateK8sBasicSideCar]
		if exists && strings.ToUpper(jo) == "TRUE" {
			podSpec.Containers = append(podSpec.Containers, v1.Container{
				Name:    jt.JobName + "-drmaa2os-sidecar",
				Image:   "drmaa/drmaa2os-basic-sidecar:latest",
				Command: []string{"/sidecar"},
				Env: []v1.EnvVar{
					{
						Name: "DRMAA2OS_POD_NAME",
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								FieldPath: "metadata.name",
							}},
					},
					{
						Name: "DRMAA2OS_POD_NAMESPACE",
						ValueFrom: &v1.EnvVarSource{
							FieldRef: &v1.ObjectFieldSelector{
								FieldPath: "metadata.namespace",
							}},
					},
				},
			})
		}
	}

	podSpec = addExtensionsAccelerators(podSpec, jt)

	podSpec.RestartPolicy = v1.RestartPolicyNever

	// settings for command etc.

	var one int32 = 1
	var zero int32 = 0

	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         jt.JobName,
			Labels:       map[string]string{"drmaa2jobsession": jobsession},
			GenerateName: "drmaa2os",
			Namespace:    namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &zero,

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

	dl, err := deadlineTime(jt)
	if err != nil {
		return nil, err
	}
	if dl != -1 {
		job.Spec.ActiveDeadlineSeconds = &dl
	}

	return addExtensions(&job, jt), nil
}

func addExtensionsAccelerators(podSpec v1.PodSpec, jt drmaa2interface.JobTemplate) v1.PodSpec {
	if jt.ExtensionList != nil {

		distribution, exists := jt.ExtensionList[extension.JobTemplateK8sDistribution]

		// AKS job using GPUs
		if exists && strings.ToLower(distribution) == "aks" {
			accelerator, set := jt.ExtensionList[extension.JobTemplateK8sAccelerator]
			if set {
				amount, gpuType := parseAccelerator(accelerator)
				if amount != "0" && gpuType != "" {
					for i := range podSpec.Containers {
						if podSpec.Containers[i].Resources.Limits == nil {
							podSpec.Containers[i].Resources.Limits = make(map[v1.ResourceName]resource.Quantity)
						}
						podSpec.Containers[i].Resources.Limits[v1.ResourceName("nvidia.com/gpu")] = resource.MustParse(amount)
					}
					if podSpec.Tolerations == nil {
						podSpec.Tolerations = make([]v1.Toleration, 0)
					}
					podSpec.Tolerations = append(podSpec.Tolerations, v1.Toleration{
						Key:      "sku",
						Operator: v1.TolerationOpEqual,
						Value:    "gpu",
						Effect:   v1.TaintEffectNoSchedule,
					})
				}
			}
		}

		// EKS job using GPUs
		if exists && strings.ToLower(distribution) == "eks" {
			accelerator, set := jt.ExtensionList[extension.JobTemplateK8sAccelerator]
			if set {
				amount, gpuType := parseAccelerator(accelerator)
				if amount != "0" && gpuType != "" {
					for i := range podSpec.Containers {
						if podSpec.Containers[i].Resources.Limits == nil {
							podSpec.Containers[i].Resources.Limits = make(map[v1.ResourceName]resource.Quantity)
						}
						podSpec.Containers[i].Resources.Limits[v1.ResourceName("nvidia.com/gpu")] = resource.MustParse(amount)
					}
				}
			}
		}

		// GKE job using GPUs
		if exists && strings.ToLower(distribution) == "gke" {
			accelerator, set := jt.ExtensionList[extension.JobTemplateK8sAccelerator]
			if set {
				amount, gpuType := parseAccelerator(accelerator)
				if amount != "0" && gpuType != "" {
					for i := range podSpec.Containers {
						if podSpec.Containers[i].Resources.Limits == nil {
							podSpec.Containers[i].Resources.Limits = make(map[v1.ResourceName]resource.Quantity)
						}
						podSpec.Containers[i].Resources.Limits[v1.ResourceName("nvidia.com/gpu")] = resource.MustParse(amount)
					}
					if podSpec.NodeSelector == nil {
						podSpec.NodeSelector = make(map[string]string)
					}
					podSpec.NodeSelector["cloud.google.com/gke-accelerator"] = gpuType
				}
			}
		}
	}
	return podSpec
}

func parseAccelerator(accelerator string) (string, string) {
	p := strings.Split(accelerator, "*")
	if len(p) != 2 {
		return "0", ""
	}
	amount, err := strconv.Atoi(p[0])
	if err != nil {
		return "0", ""
	}
	return strconv.Itoa(amount), p[1]
}
