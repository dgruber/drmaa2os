package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
			job, err := convertJob("jobsession", jt)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(BeNil())
			Ω(job.TypeMeta.Kind).Should(Equal("Job"))
			Ω(job.TypeMeta.APIVersion).Should(Equal("v1"))
			Ω(job.ObjectMeta.Name).Should(Equal("name"))
			Ω(job.Labels["drmaa2jobsession"]).Should(Equal("jobsession"))
			Ω(*job.Spec.Parallelism).Should(BeNumerically("==", 1))
			Ω(*job.Spec.Completions).Should(BeNumerically("==", 1))
		})

		It("should fail converting the JobTemplate when the JobCategory is missing", func() {
			jt.JobCategory = ""
			job, err := convertJob("", jt)
			Ω(err).ShouldNot(BeNil())
			Ω(job).Should(BeNil())
		})

		It("should fail converting the JobTemplate when the RemoteCommand is missing", func() {
			jt.RemoteCommand = ""
			job, err := convertJob("", jt)
			Ω(err).ShouldNot(BeNil())
			Ω(job).Should(BeNil())
		})

	})

})
