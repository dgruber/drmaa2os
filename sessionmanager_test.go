package drmaa2os_test

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"

	"os"

	// test with process tracker
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/singularity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const smtempdb string = "drmaa2ostest.db"

var _ = Describe("Sessionmanager", func() {

	var (
		sm drmaa2interface.SessionManager
	)

	BeforeEach(func() {
		os.Remove("drmaa2ostest")
		sm, _ = drmaa2os.NewDefaultSessionManager("drmaa2ostest")
	})

	Describe("Create and Destroy Job Session", func() {

		Context("when the Job Session does not exists", func() {
			It("should not error when creating or destroying", func() {
				js, err := sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())
				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())
			})
		})

		Context("when the Job Session already exists", func() {

			It("should error when creating", func() {
				sm.CreateJobSession("testsession", "")
				js, err := sm.CreateJobSession("testsession", "")
				Ω(err).ShouldNot(BeNil())
				Ω(js).Should(BeNil())
				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())

				js, err = sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())
				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())
			})

		})

	})

	Describe("Open a Job Session", func() {
		var js drmaa2interface.JobSession
		var err error

		BeforeEach(func() {
			js, err := sm.CreateJobSession("testsession", "")
			Ω(err).Should(BeNil())
			Ω(js).ShouldNot(BeNil())
		})

		Context("Error cases", func() {
			It("should error when the job session does not exist", func() {
				js, err = sm.OpenJobSession("doesnotexist")
				Ω(js).Should(BeNil())
				Ω(err).ShouldNot(BeNil())
			})

			It("should error when destroying a non-existing job session", func() {
				err := sm.DestroyJobSession("doesNotExist")
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("when the job session is closed", func() {

			It("should not error opening it", func() {
				js, err := sm.CreateJobSession("testsession2", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())

				err = js.Close()
				Ω(err).Should(BeNil())

				js, err = sm.OpenJobSession("testsession2")
				Ω(err).Should(BeNil())

				js.Close()
				Ω(err).Should(BeNil())

				err = sm.DestroyJobSession("testsession2")
				Ω(err).Should(BeNil())
			})

		})

		Context("when the job session is open", func() {
			It("should not error open it again (before closing)", func() {
				js, err = sm.OpenJobSession("testsession")
				Ω(err).Should(BeNil())

				js.Close()
				Ω(err).Should(BeNil())

				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())
			})
		})

		Context("when the job session is closed", func() {
			It("should be able to re-open the persistent job session", func() {
				os.Remove("drmaa2ostest")
				os.Remove("drmaa2ostestjobs")

				sm, err = drmaa2os.NewDefaultSessionManagerWithParams(
					simpletracker.SimpleTrackerInitParams{
						PersistentStorage:   true,
						PersistentStorageDB: "drmaa2ostestjobs",
					}, "drmaa2ostest")
				Ω(err).Should(BeNil())

				js, err = sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())

				err = js.Close()
				Ω(err).Should(BeNil())

				js, err = sm.OpenJobSession("testsession")
				Ω(err).Should(BeNil())

				err = js.Close()
				Ω(err).Should(BeNil())

				err := sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())

				js, err = sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())
			})
		})

	})

	Describe("Open Monitoring Session", func() {

		Context("Monitoring Session is currently not implemented", func() {
			It("should not error", func() {
				ms, err := sm.OpenMonitoringSession("")
				Ω(err).ShouldNot(BeNil())
				Ω(ms).Should(BeNil())
			})
		})

	})

	Describe("Global functions", func() {

		Describe("DRMS Name", func() {
			It("should be not an error", func() {
				name, err := sm.GetDrmsName()
				Ω(err).Should(BeNil())
				Ω(name).ShouldNot(Equal(""))
			})
		})

		Describe("DRMS Version", func() {
			It("should be not an error", func() {
				version, err := sm.GetDrmsVersion()
				Ω(err).Should(BeNil())
				Ω(version).ShouldNot(BeNil())
				Ω(version.Major).ShouldNot(Equal(""))
				Ω(version.Minor).ShouldNot(Equal(""))
			})
		})

		Describe("Support for optional functionality", func() {
			It("should be not an error", func() {
				sm.Supports(drmaa2interface.AdvanceReservation)
				sm.Supports(drmaa2interface.ReserveSlots)
				sm.Supports(drmaa2interface.Callback)
				sm.Supports(drmaa2interface.BulkJobsMaxParallel)
				sm.Supports(drmaa2interface.JtEmail)
				sm.Supports(drmaa2interface.JtStaging)
				sm.Supports(drmaa2interface.JtDeadline)
				sm.Supports(drmaa2interface.JtMaxSlots)
				sm.Supports(drmaa2interface.JtAccountingID)
				sm.Supports(drmaa2interface.RtStartNow)
				sm.Supports(drmaa2interface.RtDuration)
				sm.Supports(drmaa2interface.RtMachineOS)
				sm.Supports(drmaa2interface.RtMachineArch)
			})
		})

		Describe("Register event notification", func() {
			It("should be callable", func() {
				sm.RegisterEventNotification()
			})
		})

	})

	Describe("List sessions functionality", func() {

		Context("when they are empty", func() {
			It("should list no job and reservation sessions", func() {
				names, err := sm.GetJobSessionNames()
				Ω(err).Should(BeNil())
				Ω(names).ShouldNot(BeNil())
				Ω(len(names)).Should(BeNumerically("==", 0))
			})
		})

		Context("after adding some sessions", func() {
			It("must list all job sessions and all reservation sessions", func() {
				js, err := sm.CreateJobSession("session1", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())
				err = js.Close()
				Ω(err).Should(BeNil())

				js, err = sm.CreateJobSession("session2", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())
				err = js.Close()
				Ω(err).Should(BeNil())

				names, err := sm.GetJobSessionNames()
				Ω(err).Should(BeNil())
				Ω(len(names)).Should(BeNumerically("==", 2))
				Ω(names[0]).Should(Or(Equal("session1"), Equal("session2")))
				Ω(names[1]).Should(Or(Equal("session1"), Equal("session2")))

			})
		})

	})

	Context("ReservationSession is currently not supported", func() {

		It("should return an unsupported operation error when using a reservation session", func() {
			rs, err := sm.CreateReservationSession("reservationSession", "")
			Ω(rs).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(drmaa2os.ErrorUnsupportedOperation))

			rs, err = sm.OpenReservationSession("reservationSession")
			Ω(rs).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(drmaa2os.ErrorUnsupportedOperation))

			err = sm.DestroyReservationSession("reservationSession")
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(drmaa2os.ErrorUnsupportedOperation))

			names, err := sm.GetReservationSessionNames()
			Ω(names).Should(BeNil())
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(drmaa2os.ErrorUnsupportedOperation))
		})

	})

})
