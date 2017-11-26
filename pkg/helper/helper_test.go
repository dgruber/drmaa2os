package helper_test

import (
	. "github.com/dgruber/drmaa2os/pkg/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2os/pkg/simpletrackerfakes"
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

		It("AddArrayJobAsSingleJobs should work", func() {
			fakeTracker := simpletrackerfakes.New("testsession")
			id, err := AddArrayJobAsSingleJobs(fakeTracker, 10, 110, 2)
			Ω(err).Should(BeNil())
			jobs, errJobs := fakeTracker.ListJobs(id)
			Ω(errJobs).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 50))
		})

	})

})
