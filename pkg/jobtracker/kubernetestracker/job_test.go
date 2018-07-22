package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
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
			var fakeJobInterface fake.FakeJobs
			var fakeJob batchv1.Job

			err := jobStateChange(&fakeJobInterface, &fakeJob, "somethingwrong")
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when action is unsupported", func() {
			var fakeJobInterface fake.FakeJobs
			var fakeJob batchv1.Job

			err := jobStateChange(&fakeJobInterface, &fakeJob, "suspend")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(&fakeJobInterface, &fakeJob, "resume")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(&fakeJobInterface, &fakeJob, "hold")
			Ω(err).ShouldNot(BeNil())
			err = jobStateChange(&fakeJobInterface, &fakeJob, "release")
			Ω(err).ShouldNot(BeNil())
		})

	})

})
