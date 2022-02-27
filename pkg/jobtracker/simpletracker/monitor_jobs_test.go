package simpletracker_test

import (
	"strconv"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MonitorJobs", func() {

	It("should list all local processes", func() {
		procIDs, err := GetAllProcesses()
		Ω(err).Should(BeNil())
		Ω(len(procIDs)).Should(BeNumerically(">=", 3))
	})

	It("should return the JobInfo for a local process", func() {
		procIDs, err := GetAllProcesses()
		Ω(err).Should(BeNil())

		pid, _ := strconv.Atoi(procIDs[0])

		ji, err := GetJobInfo(int32(pid))
		Ω(err).Should(BeNil())

		Ω(ji.ID).Should(Equal(procIDs[0]))

		// should have extensions
		Ω(ji.ExtensionList).NotTo(BeNil())
	})

})
