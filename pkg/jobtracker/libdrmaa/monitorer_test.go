package libdrmaa

import (
	"log"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Monitorer", func() {

	createTracker := func(standard bool) *DRMAATracker {
		if standard {
			log.Println("using standard tracker")
			standardTracker, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			return standardTracker
		}
		log.Println("using tracker with persistent job storage")
		params := LibDRMAASessionParams{
			ContactString:           "",
			UsePersistentJobStorage: true,
			DBFilePath:              getTempFile(),
		}
		trackerWithParams, err := NewDRMAATrackerWithParams(params)
		Expect(err).To(BeNil())
		return trackerWithParams
	}

	Context("Grid Engine JobTracker", func() {

		It("should should implement the Monitorer interface", func() {
			tracker := createTracker(true)
			defer tracker.DestroySession()
			stracker := jobtracker.JobTracker(tracker)
			_, hasInterface := stracker.(jobtracker.Monitorer)
			Expect(hasInterface).To(BeTrue())
		})

		It("should should open and close a monitoring session", func() {
			tracker := createTracker(true)
			defer tracker.DestroySession()
			stracker := jobtracker.JobTracker(tracker)
			m, hasInterface := stracker.(jobtracker.Monitorer)
			Expect(hasInterface).To(BeTrue())
			err := m.OpenMonitoringSession("testX")
			Expect(err).To(BeNil())
			err = m.CloseMonitoringSession("testX")
			Expect(err).To(BeNil())
		})

		It("should return the local test machine", func() {
			tracker := createTracker(true)
			defer tracker.DestroySession()
			machines, err := tracker.GetAllMachines(nil)
			Expect(err).To(BeNil())
			Expect(len(machines)).To(BeNumerically("==", 1))

			// now with filter for all machines
			filter := []string{"ThisM", "achineDoes", "notExist"}
			machines, err = tracker.GetAllMachines(filter)
			Expect(err).To(BeNil())
			Expect(len(machines)).To(BeNumerically("==", 0))
		})

		It("should return monitoring jobs", func() {
			stracker := createTracker(true)
			defer stracker.DestroySession()
			_, err := stracker.GetAllJobIDs(nil)
			Expect(err).To(BeNil())
		})

	})

})
