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
				JobName:       "name",
				RemoteCommand: "command",
				Args:          []string{"arg1", "arg2"},
				JobCategory:   "category",
			}
		})

		It("should create a container spec out of the JobTemplate", func() {
			c, err := newContainers(jt)
			Ω(err).Should(BeNil())
			Ω(c).ShouldNot(BeNil())
			Ω(len(c)).Should(BeNumerically("==", 1))
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

})
