package kubernetestracker

import (
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MonitorerTest", func() {

	Context("Basisc Operations", func() {

		var kt jobtracker.Monitorer

		BeforeEach(func() {
			var err error
			kt, err = New("monitorertest", "default", nil)
			Ω(err).Should(BeNil())

		})

		It("should be possible to list all jobs", func() {
			jobs, err := kt.GetAllJobIDs(nil)
			Ω(err).Should(BeNil())
			Ω(jobs).ShouldNot(BeNil())
			if len(jobs) > 0 {
				jobInfo, err := kt.JobInfoFromMonitor(jobs[0])
				Ω(err).Should(BeNil())
				Ω(jobInfo.ID).Should(Equal(jobs[0]))
			}
		})

		It("should open and close a monitoring session without error", func() {
			err := kt.OpenMonitoringSession("")
			Ω(err).Should(BeNil())
			err = kt.CloseMonitoringSession("")
			Ω(err).Should(BeNil())
		})

		It("should open and close a monitoring session without error", func() {
			queues, err := kt.GetAllQueueNames(nil)
			Ω(err).Should(BeNil())
			Ω(queues).ShouldNot(BeNil())
			// hopefully there is a default namespace
			Ω(queues).Should(ContainElement(ContainSubstring("default")))
		})

		It("should list all machines", func() {
			machines, err := kt.GetAllMachines(nil)
			Ω(err).Should(BeNil())
			Ω(machines).ShouldNot(BeNil())
		})

	})

})
