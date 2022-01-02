package drmaa2os_test

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	"os"

	// test with process tracker
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/singularity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MonitoringSession", func() {

	var (
		sm drmaa2interface.SessionManager
	)

	BeforeEach(func() {
		os.Remove("drmaa2ostest")
		sm, _ = drmaa2os.NewDefaultSessionManager("drmaa2ostest")
	})

	Describe("Monitoring Session", func() {

		Context("Monitoring Session is implemented for process backend", func() {

			It("should create a usable monitoring session", func() {
				ms, err := sm.OpenMonitoringSession("test")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				machines, err := ms.GetAllMachines(nil)
				Ω(err).Should(BeNil())
				Ω(machines).ShouldNot(BeNil())

				jobs, err := ms.GetAllJobs(drmaa2interface.CreateJobInfo())
				Ω(err).Should(BeNil())
				Ω(jobs).ShouldNot(BeNil())

				queues, err := ms.GetAllQueues(nil)
				Ω(err).Should(BeNil())
				Ω(queues).ShouldNot(BeNil())

				// reservations are not yet supported
				reservations, err := ms.GetAllReservations()
				Ω(err).ShouldNot(BeNil())
				Ω(reservations).Should(BeNil())

				err = ms.CloseMonitoringSession()
				Ω(err).Should(BeNil())
			})

		})

		Context("Monitoring Session jobs", func() {

			It("should fail to manipulate the jobs as they are read only (for now)", func() {

				ms, err := sm.OpenMonitoringSession("test")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				jobs, err := ms.GetAllJobs(drmaa2interface.CreateJobInfo())
				Ω(err).Should(BeNil())

				// maninulation must fail
				Ω(jobs[0].Suspend()).ShouldNot(BeNil())
				Ω(jobs[0].Resume()).ShouldNot(BeNil())
				Ω(jobs[0].Hold()).ShouldNot(BeNil())
				Ω(jobs[0].Terminate()).ShouldNot(BeNil())
				Ω(jobs[0].Release()).ShouldNot(BeNil())

				// reaping a monitoring job must fail
				Ω(jobs[0].Reap()).ShouldNot(BeNil())

			})

			It("should get the state of a job", func() {

				ms, err := sm.OpenMonitoringSession("test")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				jobs, err := ms.GetAllJobs(drmaa2interface.CreateJobInfo())
				Ω(err).Should(BeNil())

				Ω(jobs[0].GetState().String()).Should(Equal(drmaa2interface.Running.String()))

				jobInfo, err := jobs[0].GetJobInfo()
				fmt.Printf("%v", jobInfo)
				Ω(err).Should(BeNil())
				Ω(jobInfo.ID).ShouldNot(Equal(""))

			})

		})

	})

})
