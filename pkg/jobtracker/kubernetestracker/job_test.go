package kubernetestracker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes/typed/batch/v1/fake"
)

var _ = Describe("Job", func() {

	Context("Error situations for jobStateChange", func() {

		It("should error when given job is nil", func() {
			err := jobStateChange(nil, nil, "terminate")
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when action is undefined", func() {
			var fakeJobInterface fake.FakeBatchV1
			var fakeJob batchv1.Job

			err := jobStateChange(fakeJobInterface.Jobs("default"),
				&fakeJob, "somethingwrong")
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when action is unsupported", func() {
			var fakeJobInterface fake.FakeBatchV1
			var fakeJob batchv1.Job

			err := jobStateChange(fakeJobInterface.Jobs("default"),
				&fakeJob, "suspend")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(fakeJobInterface.Jobs("default"),
				&fakeJob, "resume")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(fakeJobInterface.Jobs("default"),
				&fakeJob, "hold")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(fakeJobInterface.Jobs("default"),
				&fakeJob, "release")
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when job is not found", func() {
			ji, jc, err := getJobInterfaceAndJob(nil, "x", "default")
			Ω(err).ShouldNot(BeNil())
			Ω(ji).Should(BeNil())
			Ω(jc).Should(BeNil())
		})

	})

	Context("Standard error cases", func() {
		It("should error when given job is nil", func() {
			err := deleteJob(nil, nil)
			Ω(err).ShouldNot(BeNil())
		})
	})

})
