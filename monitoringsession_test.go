package drmaa2os_test

import (
	"fmt"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	"os"

	// test with process tracker
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
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

		Context("Filter", func() {

			It("should filter for a job with a specific ID", func() {
				ms, err := sm.OpenMonitoringSession("testFilter")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				jobs, err := ms.GetAllJobs(drmaa2interface.CreateJobInfo())
				Ω(err).Should(BeNil())
				Ω(jobs).ShouldNot(BeNil())
				Ω(len(jobs)).Should(BeNumerically(">=", 0))

				filter := drmaa2interface.CreateJobInfo()
				filter.ID = jobs[0].GetID()

				// hopefully the process is still running :)
				filteredJobs, err := ms.GetAllJobs(filter)
				Ω(err).Should(BeNil())
				Ω(filteredJobs).ShouldNot(BeNil())
				Ω(len(filteredJobs)).Should(BeNumerically("==", 1))
				Ω(filteredJobs[0].GetID()).Should(Equal(filter.ID))
			})

			It("should filter for the local machine", func() {
				ms, err := sm.OpenMonitoringSession("test")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				machines, err := ms.GetAllMachines(nil)
				Ω(err).Should(BeNil())
				Ω(machines).ShouldNot(BeNil())
				Ω(len(machines)).Should(BeNumerically("==", 1))

				filteredMachines, err := ms.GetAllMachines([]string{"XXX"})
				Ω(err).Should(BeNil())
				Ω(len(filteredMachines)).Should(BeNumerically("==", 0))

				filteredMachines2, err := ms.GetAllMachines([]string{machines[0].Name})
				Ω(err).Should(BeNil())
				Ω(len(filteredMachines2)).Should(BeNumerically("==", 1))
			})

			It("should filter for queues", func() {
				ms, err := sm.OpenMonitoringSession("test")
				Ω(err).Should(BeNil())
				Ω(ms).ShouldNot(BeNil())

				// there are no queues for processes
				queues, err := ms.GetAllQueues(nil)
				Ω(err).Should(BeNil())
				Ω(queues).ShouldNot(BeNil())

				queues, err = ms.GetAllQueues([]string{})
				Ω(err).Should(BeNil())
				Ω(queues).ShouldNot(BeNil())

				queues, err = ms.GetAllQueues([]string{"X"})
				Ω(err).Should(BeNil())
				Ω(queues).ShouldNot(BeNil())

			})

		})

	})

})
