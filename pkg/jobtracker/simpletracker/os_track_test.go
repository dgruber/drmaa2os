package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"os/exec"
)

var _ = Describe("OsProcessSupervisor", func() {

	It("should detect a successful run", func() {
		jobid := "1"
		ch := make(chan JobEvent, 1)
		cmd := exec.Command("true")
		err := cmd.Start()
		Ω(err).Should(BeNil())
		TrackProcess(cmd, jobid, ch, 0, nil)

		var je JobEvent
		Ω(ch).Should(Receive(&je))

		Ω(je.JobID).Should(Equal("1"))
		Ω(je.JobState).Should(Equal(drmaa2interface.Done))

		// check job info
		Ω(je.JobInfo.ID).Should(Equal("1"))
		Ω(je.JobInfo.ExitStatus).Should(BeNumerically("==", 0))
	})

	It("should detect an execution failure", func() {
		jobid := "1"
		ch := make(chan JobEvent, 1)
		cmd := exec.Command("false")
		err := cmd.Start()
		Ω(err).Should(BeNil())
		TrackProcess(cmd, jobid, ch, 0, nil)

		var je JobEvent
		Ω(ch).Should(Receive(&je))

		Ω(je.JobID).Should(Equal("1"))
		Ω(je.JobState).Should(Equal(drmaa2interface.Failed))

		// check job info
		Ω(je.JobInfo.ID).Should(Equal("1"))
		Ω(je.JobInfo.ExitStatus).Should(BeNumerically("==", 1))
	})

	It("should return a job info object", func() {

	})

})
