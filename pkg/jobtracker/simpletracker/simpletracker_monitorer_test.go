package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MonitorHost", func() {

	Context("Monitorer interface implementation", func() {

		jt := New("testsession")

		err := jt.OpenMonitoringSession("testmonitoringsession")
		Ω(err).Should(BeNil())

		It("should list all process IDs", func() {
			ids, err := jt.GetAllJobIDs(nil)
			Ω(err).Should(BeNil())

			Ω(len(ids)).Should(BeNumerically(">=", 3))
		})

		It("should return JobInfo for a process ID", func() {
			// check with init process
			ji, err := jt.JobInfoFromMonitor("1")
			Ω(err).Should(BeNil())

			Ω(ji.ID).Should(Equal("1"))
		})

		It("should return an empty queue list as we don't have queues in OSs", func() {
			ids, err := jt.GetAllQueueNames(nil)
			Ω(err).Should(BeNil())
			Ω(len(ids)).Should(BeNumerically(">=", 0))
		})

		It("should list the local machines", func() {
			machines, err := jt.GetAllMachines(nil)
			Ω(err).Should(BeNil())
			Ω(len(machines)).Should(BeNumerically(">=", 1))
			Ω(machines[0].Name).ShouldNot(Equal(""))
			Ω(machines[0].Available).Should(BeTrue())
			Ω(machines[0].Load).ShouldNot(BeZero())
			Ω(machines[0].PhysicalMemory).ShouldNot(BeZero())
			Ω(machines[0].VirtualMemory).ShouldNot(BeZero())
		})

		It("should open and cluster a monitoring session", func() {
			jt := New("testopensession")
			err := jt.OpenMonitoringSession("ms")
			Ω(err).Should(BeNil())
			err = jt.CloseMonitoringSession("ms")
			Ω(err).Should(BeNil())
		})

		It("should filter machines", func() {
			machines, err := jt.GetAllMachines([]string{"NoTeXistingMachine"})
			Ω(err).Should(BeNil())
			Ω(len(machines)).Should(BeNumerically(">=", 0))
		})

		It("should filter machines", func() {
			machines, err := jt.GetAllMachines([]string{"NoTeXistingMachine"})
			Ω(err).Should(BeNil())
			Ω(len(machines)).Should(BeNumerically(">=", 0))
		})

	})

})
