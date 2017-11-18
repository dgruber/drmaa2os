package kubernetestracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KubernetesTracker", func() {

	Context("Basic interface test", func() {

		jt := drmaa2interface.JobTemplate{
			JobName:       "name",
			RemoteCommand: "command",
			JobCategory:   "image",
			Args:          []string{"123"},
		}

		It("should be possible to create a KubernetesTracker", func() {
			kt, err := New()
			Ω(err).Should(BeNil())
			Ω(kt).ShouldNot(BeNil())
		})

		It("should be possible to AddJob()", func() {
			kt, err := New()
			Ω(err).Should(BeNil())
			Ω(kt).ShouldNot(BeNil())

			jobid, err := kt.AddJob(jt)
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

	})

})
