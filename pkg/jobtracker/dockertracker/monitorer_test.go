package dockertracker_test

import (
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Monitorer Interface of DockerTracker", func() {

	Context("basic functionality", func() {

		var tracker jobtracker.Monitorer

		BeforeEach(func() {
			tracker, _ = New("")
		})

		It("should open and close a monitoring session without issues", func() {
			err := tracker.OpenMonitoringSession("test")
			Ω(err).Should(BeNil())
			err = tracker.CloseMonitoringSession("test")
			Ω(err).Should(BeNil())
		})

		It("should return all containers without filter", func() {
			containers, err := tracker.GetAllJobIDs(nil)
			Ω(err).Should(BeNil())
			Ω(containers).ShouldNot(BeNil())
			Ω(len(containers)).Should(BeNumerically(">=", 1))

			jobinfo, err := tracker.JobInfoFromMonitor(containers[0])
			Ω(err).Should(BeNil())
			Ω(jobinfo.ID).Should(Equal(containers[0]))
		})

		It("should return the docker host", func() {
			dockerhost, err := tracker.GetAllMachines(nil)
			Ω(err).Should(BeNil())
			Ω(dockerhost).ShouldNot(BeNil())
			Ω(len(dockerhost)).Should(BeNumerically("==", 1))
		})

		It("should return the empty queue names", func() {
			queues, err := tracker.GetAllQueueNames(nil)
			Ω(err).Should(BeNil())
			Ω(queues).ShouldNot(BeNil())
			Ω(len(queues)).Should(BeNumerically("==", 0))
		})

	})

})
