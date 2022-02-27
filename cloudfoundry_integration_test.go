// +build cf_integration

package drmaa2os_test

import (
	. "github.com/dgruber/drmaa2os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"os"
)

var _ = Describe("Cloud Foundry integration", func() {
	var sm drmaa2interface.SessionManager
	var js drmaa2interface.JobSession
	var jt drmaa2interface.JobTemplate

	BeforeEach(func() {
		os.Remove("cintegration_tmp.db")
		var err error

		// test expects that these environment variables are set
		Ω(os.Getenv("CF_INSTANCE_GUID")).ShouldNot(Equal(""))
		Ω(os.Getenv("CF_API")).ShouldNot(Equal(""))
		Ω(os.Getenv("CF_USER")).ShouldNot(Equal(""))
		Ω(os.Getenv("CF_PASSWORD")).ShouldNot(Equal(""))

		sm, err = NewCloudFoundrySessionManager(os.Getenv("CF_API"), os.Getenv("CF_USER"), os.Getenv("CF_PASSWORD"), "cintegration_tmp.db")
		Ω(err).Should(BeNil())
		js, err = sm.CreateJobSession("testsession", "")
		Ω(err).Should(BeNil())

		jt = drmaa2interface.JobTemplate{
			RemoteCommand: "/bin/sleep",
			Args:          []string{"13"},
			JobCategory:   os.Getenv("CF_INSTANCE_GUID"),
		}
	})

	It("submits a task", func() {
		job, err := js.RunJob(jt)
		Ω(err).Should(BeNil())
		Ω(job).ShouldNot(BeNil())

	})

	It("submits a task and waits for it", func() {
		jt.Args = []string{"0"}
		job, err := js.RunJob(jt)
		Ω(err).Should(BeNil())
		Ω(job).ShouldNot(BeNil())
		Ω(job.WaitTerminated(drmaa2interface.InfiniteTime)).Should(BeNil())
		Ω(job.GetState()).Should(Equal(drmaa2interface.Done))
		ji, err := job.GetJobInfo()
		Ω(err).Should(BeNil())
		Ω(ji.ExitStatus).Should(BeNumerically("==", 0))
	})

	It("submits a failing task and waits for it", func() {
		jt.RemoteCommand = "notacommand"
		job, err := js.RunJob(jt)
		Ω(err).Should(BeNil())
		Ω(job).ShouldNot(BeNil())
		Ω(job.WaitTerminated(drmaa2interface.InfiniteTime)).Should(BeNil())
		Ω(job.GetState()).Should(Equal(drmaa2interface.Failed))
	})

	It("submits a task array and waits for it", func() {
		jt.Args = []string{"0"}
		aj, err := js.RunBulkJobs(jt, 1, 3, 1, -1)
		Ω(err).Should(BeNil())
		Ω(aj).ShouldNot(BeNil())
		Ω(len(aj.GetJobs())).Should(BeNumerically("==", 3))
		for _, job := range aj.GetJobs() {
			Ω(err).Should(BeNil())
			Ω(job.WaitTerminated(drmaa2interface.InfiniteTime)).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Done))
			ji, err := job.GetJobInfo()
			Ω(err).Should(BeNil())
			Ω(ji.ExitStatus).Should(BeNumerically("==", 0))
		}
	})

})
