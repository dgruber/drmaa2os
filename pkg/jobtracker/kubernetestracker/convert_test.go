package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/api/batch/v1"
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
		})

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
			Ω(err.Error()).Should(Equal("JobCategory (image name) not set in JobTemplate"))
		})

		It("should convert the JobTemplate into a Job", func() {
			job, err := convertJob(jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())

			Ω(job.TypeMeta.Kind).Should(Equal("Job"))
			Ω(job.TypeMeta.APIVersion).Should(Equal("v1"))

			Ω(job.ObjectMeta.Name).Should(Equal("name"))

			Ω(*job.Spec.Parallelism).Should(BeNumerically("==", 1))
			Ω(*job.Spec.Completions).Should(BeNumerically("==", 1))
		})

	})

	Context("JobStatus conversion", func() {
		var status batchv1.JobStatus

		BeforeEach(func() {
			status = batchv1.JobStatus{
				Active:    0,
				Succeeded: 0,
				Failed:    0,
			}
		})

		It("should convert nil to Undetermined state", func() {
			Ω(convertJobStatus2JobState(nil)).Should(Equal(drmaa2interface.Undetermined))
		})

		It("should convert active to Running state", func() {
			status.Active = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Running))
		})

		It("should convert failed to Failed state", func() {
			status.Failed = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Failed))
		})

		It("should convert succeeded to Done state", func() {
			status.Succeeded = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Done))
		})

		It("should convert unset states to Undetermined state", func() {
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Undetermined))
		})

	})

})
