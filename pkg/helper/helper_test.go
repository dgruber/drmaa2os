package helper_test

import (
	. "github.com/dgruber/drmaa2os/pkg/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletrackerfakes"
)

var _ = Describe("Helper", func() {

	Context("Array Job ID convert functions", func() {

		It("should generate and resolve an array job ID into job IDs", func() {
			guids := []string{"1", "2", "3"}

			id := Guids2ArrayJobID(guids)
			guidsOut, err := ArrayJobID2GUIDs(id)

			Ω(err).Should(BeNil())
			Ω(guidsOut).Should(BeEquivalentTo(guids))
		})

	})

	Context("Create array job out with single job submissions", func() {

		jt := drmaa2interface.JobTemplate{RemoteCommand: "test"}

		It("AddArrayJobAsSingleJobs should work", func() {
			fakeTracker := simpletrackerfakes.New("testsession")
			_, err := AddArrayJobAsSingleJobs(jt, fakeTracker, 11, 110, 2)
			Ω(err).Should(BeNil())
			jobs, errJobs := fakeTracker.ListJobs()
			Ω(errJobs).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 50))
		})

	})

	Context("Error cases", func() {

		It("should return nothing when array job id is not parsable", func() {
			ajid, err := ArrayJobID2GUIDs("_")
			Ω(err).ShouldNot(BeNil())
			Ω(ajid).Should(BeEmpty())
		})
	})

})
