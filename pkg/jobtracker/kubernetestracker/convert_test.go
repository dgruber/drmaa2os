package kubernetestracker

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Convert", func() {

	Context("job conversion", func() {

		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				JobName:          "name",
				RemoteCommand:    "command",
				Args:             []string{"arg1", "arg2"},
				JobCategory:      "category",
				WorkingDirectory: "/workingdirectory",
				JobEnvironment: map[string]string{
					"ENV1": "CONTENT1",
					"ENV2": "CONTENT2",
				},
			}
		})

		It("should create a container spec out of the JobTemplate", func() {
			c, err := newContainers(jt)
			Ω(err).Should(BeNil())
			Ω(c).ShouldNot(BeNil())
			Ω(len(c)).Should(BeNumerically("==", 1))

			c0 := c[0]
			Ω(c0.Name).Should(Equal(jt.JobName))
			Ω(c0.Image).Should(Equal(jt.JobCategory))
			Ω(c0.Command[0]).Should(Equal(jt.RemoteCommand))
			Ω(c0.Args).Should(BeEquivalentTo(jt.Args))
			Ω(c0.WorkingDir).Should(Equal(jt.WorkingDirectory))
			Ω(len(c0.Env)).Should(BeNumerically("==", 2))
			Ω(c0.Env[0].Name).Should(Or(Equal("ENV1"), Equal("ENV2")))
			Ω(c0.Env[1].Name).Should(Or(Equal("ENV1"), Equal("ENV2")))
			Ω(c0.Env[0].Value).Should(Or(Equal("CONTENT1"), Equal("CONTENT2")))
			Ω(c0.Env[1].Value).Should(Or(Equal("CONTENT1"), Equal("CONTENT2")))
		})

		It("should convert the JobTemplate into a Job", func() {
			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			Ω(job.TypeMeta.Kind).Should(Equal("Job"))
			Ω(job.TypeMeta.APIVersion).Should(Equal("v1"))
			Ω(job.ObjectMeta.Name).Should(Equal("name"))
			Ω(job.Labels["drmaa2jobsession"]).Should(Equal("jobsession"))
			Ω(*job.Spec.Parallelism).Should(BeNumerically("==", 1))
			Ω(*job.Spec.Completions).Should(BeNumerically("==", 1))
		})

		It("should convert the deadline from time to int", func() {
			jt.DeadlineTime = time.Now().Add(time.Second * 10)
			deadline, err := deadlineTime(jt)
			Ω(err).Should(BeNil())
			Ω(*deadline).Should(BeNumerically("<=", 10))
		})

		It("should add a label to the job object when requested as extension", func() {
			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			jt.ExtensionList = map[string]string{"labels": "label1=foo,label2=bar,drmaa2jobsession=UI"}
			job = addExtensions(job, jt)
			Ω(job.Labels["label1"]).Should(Equal("foo"))
			Ω(job.Labels["label2"]).Should(Equal("bar"))
			// not allowed to override job session
			Ω(job.Labels["drmaa2jobsession"]).Should(Equal("jobsession"))
		})

		It("should select a scheduler when requested as extension", func() {
			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			jt.ExtensionList = map[string]string{"scheduler": "poseidon"}
			job = addExtensions(job, jt)
			Ω(job.Spec.Template.Spec.SchedulerName).Should(Equal("poseidon"))
		})

		It("should run privileged when 'privileged' is set as JobTemplate extension", func() {
			jt.ExtensionList = map[string]string{"privileged": "true"}

			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())

			Ω(job.Spec.Template.Spec.Containers[0]).ShouldNot(BeNil())
			Ω(job.Spec.Template.Spec.Containers[0].SecurityContext).ShouldNot(BeNil())
			Ω(*job.Spec.Template.Spec.Containers[0].SecurityContext.Privileged).To(BeTrue())
		})

		It("should not run privileged when 'privileged' is set as JobTemplate extension", func() {
			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			Ω(job.Spec.Template.Spec.Containers[0].SecurityContext).To(BeNil())
		})

		It("should add env variables from configmaps and secrets", func() {
			jt.ExtensionList = map[string]string{
				"env-from-secrets":    "secretname1:secretname2",
				"env-from-configmaps": "configmap1:configmap2",
			}
			job, err := convertJob("jobsession", "default", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			Ω(job.Spec.Template.Spec.Containers[0].EnvFrom).NotTo(BeNil())
			Ω(job.Spec.Template.Spec.Containers[0].EnvFrom[0].SecretRef.Name).To(Equal("secretname1"))
			Ω(job.Spec.Template.Spec.Containers[0].EnvFrom[1].SecretRef.Name).To(Equal("secretname2"))
			Ω(job.Spec.Template.Spec.Containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("configmap1"))
			Ω(job.Spec.Template.Spec.Containers[0].EnvFrom[3].ConfigMapRef.Name).To(Equal("configmap2"))
		})

		Context("error cases", func() {
			It("should error when the RemoteCommand is not set in the JobTemplate", func() {
				jt.RemoteCommand = ""
				c, err := newContainers(jt)
				Ω(c).Should(BeNil())
				Ω(err).ShouldNot(BeNil())
				Ω(err.Error()).Should(Equal("RemoteCommand not set in JobTemplate"))
			})

			It("should error when the JobCategory is not set in the JobTemplate", func() {
				jt.JobCategory = ""
				c, err := newContainers(jt)
				Ω(c).Should(BeNil())
				Ω(err).ShouldNot(BeNil())
				Ω(err.Error()).Should(Equal("JobCategory (container image name) not set in JobTemplate"))
			})

			It("should fail converting the JobTemplate when the JobCategory is missing", func() {
				jt.JobCategory = ""
				job, err := convertJob("", "default", jt)
				Ω(err).ShouldNot(BeNil())
				Ω(job).Should(BeNil())
			})

			It("should fail converting the JobTemplate when the RemoteCommand is missing", func() {
				jt.RemoteCommand = ""
				job, err := convertJob("", "default", jt)
				Ω(err).ShouldNot(BeNil())
				Ω(job).Should(BeNil())
			})

			It("should fail converting the JobTemplate when DeadlineTime is in the past", func() {
				jt.DeadlineTime = time.Now().Add(time.Second * -1)
				job, err := convertJob("", "default", jt)
				Ω(err).ShouldNot(BeNil())
				Ω(job).Should(BeNil())
			})
		})

		Context("File staging: Local mounts (hostpath)", func() {

			var jt drmaa2interface.JobTemplate

			BeforeEach(func() {
				jt = drmaa2interface.JobTemplate{
					JobCategory:   "busybox:latest",
					JobName:       "name",
					RemoteCommand: "/entrypoint.sh",
				}
				jt.StageInFiles = map[string]string{
					"/root":             "hostpath:/",
					"/usr/local/nvidia": "hostpath:/home/kubernetes/bin/nvidia",
				}
				jt.ExtensionList = map[string]string{
					"privileged": "true",
				}
			})

			It("should create new volumes", func() {
				v, err := newVolumes(jt)
				Ω(err).Should(BeNil())
				Ω(len(v)).Should(BeNumerically("==", 2))
				Ω(v[0].HostPath).ShouldNot(BeNil())
				Ω(v[0].HostPath.Path).Should(Or(Equal("/"), Equal("/home/kubernetes/bin/nvidia")))
				Ω(v[1].HostPath).ShouldNot(BeNil())
				Ω(v[1].HostPath.Path).Should(Or(Equal("/"), Equal("/home/kubernetes/bin/nvidia")))
			})

			It("should convert the JobTemplate to a job spec with hostpath volumes and mounts", func() {
				job, err := convertJob("jobsession", "default", jt)
				Ω(err).Should(BeNil())
				Ω(job).ShouldNot(BeNil())
				Ω(job.Spec.Template.Spec.Volumes).ShouldNot(BeNil())
				Ω(len(job.Spec.Template.Spec.Volumes)).Should(BeNumerically("==", 2))
				Ω(len(job.Spec.Template.Spec.Containers[0].VolumeMounts)).Should(BeNumerically("==", 2))
				Ω(job.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).Should(Or(Equal("/root"), Equal("/usr/local/nvidia")))
				Ω(job.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).Should(Or(Equal("/root"), Equal("/usr/local/nvidia")))
			})

		})

		Context("File staging: ConfigMaps and Secrets", func() {

			var jt drmaa2interface.JobTemplate

			BeforeEach(func() {
				jt = drmaa2interface.JobTemplate{
					JobName: "name",
				}
				jt.StageInFiles = map[string]string{
					"/my/secret.txt":            "secret-data:c2VjcmV0Cg==",
					"/my/configmap.txt":         "configmap-data:c2VjcmV0Cg==",
					"/my/othersecret.txt":       "secret-data:c2VjcmV0Mgo=",
					"/existing/input.txt":       "configmap:existingConfigMapName",
					"/existing/secretInput.txt": "secret:existingSecretName",
					"/existing/pvc":             "pvc:existingPVCName",
					"/nfs":                      "nfs:server:/share/directory",
				}
			})

			It("should create appropriate names for configmaps, secrets, and volumes", func() {
				v := volumeName("job123", "/mount/file/here.txt", "secret")
				Ω(v).Should(ContainSubstring("job123-"))
				Ω(v).Should(ContainSubstring("-secret-volume"))

				v2 := volumeName("job123", "/mount/file/here2.txt", "configmap")
				Ω(v2).Should(ContainSubstring("job123-"))
				Ω(v2).Should(ContainSubstring("-configmap-volume"))

				Ω(v).ShouldNot(Equal(v2))

				cm1 := configMapName("job123", "/mount/file/bla")
				cm2 := configMapName("job123", "/mount/file/bla2")
				Ω(cm1).ShouldNot(Equal(cm2))

				Ω(cm1).Should(ContainSubstring("job123-"))
				Ω(cm2).Should(ContainSubstring("job123-"))

				s1 := secretName("job123", "/mount/file/bla")
				s2 := secretName("job123", "/mount/file/bla2")
				Ω(s1).ShouldNot(Equal(s2))

				Ω(s1).Should(ContainSubstring("job123-"))
				Ω(s2).Should(ContainSubstring("job123-"))
			})

			It("should create new volumes", func() {
				v, err := newVolumes(jt)
				Ω(err).Should(BeNil())
				Ω(len(v)).Should(BeNumerically("==", 7))
				Ω(v[0].Name).ShouldNot(Equal(""))
				Ω(v[1].Name).ShouldNot(Equal(""))
				Ω(v[2].Name).ShouldNot(Equal(""))
				Ω(v[3].Name).ShouldNot(Equal(""))
				Ω(v[4].Name).ShouldNot(Equal(""))
				Ω(v[5].Name).ShouldNot(Equal(""))
				Ω(v[6].Name).ShouldNot(Equal(""))
				Ω(v[0].VolumeSource).ShouldNot(BeNil())
				Ω(v[1].VolumeSource).ShouldNot(BeNil())
				Ω(v[2].VolumeSource).ShouldNot(BeNil())
				Ω(v[3].VolumeSource).ShouldNot(BeNil())
				Ω(v[4].VolumeSource).ShouldNot(BeNil())
				Ω(v[5].VolumeSource).ShouldNot(BeNil())
				Ω(v[6].VolumeSource).ShouldNot(BeNil())
				// contains the unique job id
				Ω(v[0].Name).Should(ContainSubstring("name"))
				Ω(v[1].Name).Should(ContainSubstring("name"))
				Ω(v[2].Name).Should(ContainSubstring("name"))
				Ω(v[3].Name).Should(ContainSubstring("name"))
				Ω(v[4].Name).Should(ContainSubstring("name"))
				Ω(v[5].Name).Should(ContainSubstring("name"))
				Ω(v[6].Name).Should(ContainSubstring("name"))

				pvc := 0
				nfs := 0
				for _, vol := range v {
					if strings.Contains(vol.Name, "pvc") {
						Ω(vol.PersistentVolumeClaim.ClaimName).Should(Equal("existingPVCName"))
						pvc++
					}
					if strings.Contains(vol.Name, "nfs") {
						Ω(vol.NFS).ShouldNot(BeNil())
						Ω(vol.NFS.Server).Should(Equal("server"))
						Ω(vol.NFS.Path).Should(Equal("/share/directory"))
						nfs++
					}
				}
				Ω(pvc).Should(BeNumerically("==", 1))
				Ω(nfs).Should(BeNumerically("==", 1))
			})

		})

	})

})
