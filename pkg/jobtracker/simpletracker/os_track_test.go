package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"os/exec"
	"time"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("OsProcessSupervisor", func() {

	It("should detect a successful run", func() {
		jobid := "1"
		ch := make(chan JobEvent, 1)
		cmd := exec.Command("true")
		err := cmd.Start()
		Ω(err).Should(BeNil())
		TrackProcess(cmd, nil, jobid, time.Now(), ch, 0, nil)

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
		TrackProcess(cmd, nil, jobid, time.Now(), ch, 0, nil)

		var je JobEvent
		Ω(ch).Should(Receive(&je))

		Ω(je.JobID).Should(Equal("1"))
		Ω(je.JobState).Should(Equal(drmaa2interface.Failed))

		// check job info
		Ω(je.JobInfo.ID).Should(Equal("1"))
		Ω(je.JobInfo.ExitStatus).Should(BeNumerically("==", 1))

		// check for extensions
		Ω(je.JobInfo.ExtensionList).ShouldNot(BeNil())
		_, exists := je.JobInfo.ExtensionList["system_time_ms"]
		Ω(exists).Should(BeTrue())
		_, exists = je.JobInfo.ExtensionList["user_time_ms"]
		Ω(exists).Should(BeTrue())
	})

})
