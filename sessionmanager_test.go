package drmaa2os

import (
	"code.cloudfoundry.org/lager"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/storage/boltstore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

const smtempdb string = "drmaa2ostest.db"

func createSessionManager() drmaa2interface.SessionManager {
	os.Remove(smtempdb)
	s := boltstore.NewBoltStore(smtempdb)
	s.Init()
	l := lager.NewLogger("drmaa2ostest")
	l.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
	return &SessionManager{store: s, log: l}
}

var _ = Describe("Sessionmanager", func() {

	Describe("Create and Destroy Job Session", func() {

		var (
			sm drmaa2interface.SessionManager
			//js drmaa2interface.JobSession
		)

		BeforeEach(func() {
			os.Remove("drmaa2ostest")
			sm, _ = NewDefaultSessionManager("drmaa2ostest")
		})

		Context("when the Job Session does not exist", func() {
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

		var (
			sm drmaa2interface.SessionManager
		)

		BeforeEach(func() {
			sm = createSessionManager()
		})

		Context("when the job session is closed", func() {
			It("should not error", func() {
				js, err := sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())

				err = js.Close()
				Ω(err).Should(BeNil())

				js, err = sm.OpenJobSession("testsession")
				Ω(err).Should(BeNil())

				js.Close()
				Ω(err).Should(BeNil())

				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())
			})
		})

		Context("when the job session is open", func() {
			It("should not error", func() {
				js, err := sm.CreateJobSession("testsession", "")
				Ω(err).Should(BeNil())
				Ω(js).ShouldNot(BeNil())

				js, err = sm.OpenJobSession("testsession")
				Ω(err).Should(BeNil())

				js.Close()
				Ω(err).Should(BeNil())

				err = sm.DestroyJobSession("testsession")
				Ω(err).Should(BeNil())
			})
		})

	})

	Describe("Open Monitoring Session", func() {
		/*
			var (
				sm drmaa2interface.SessionManager
			)

			BeforeEach(func() {
				sm = createSessionManager()
			})

			Context("when the Monitoring Session does not exist", func() {
				It("should not error", func() {
				ms, err := sm.OpenMonitoringSession("")
					Ω(err).Should(BeNil())
					Ω(ms).ShouldNot(BeNil())
					err = ms.CloseMonitoringSession()
					Ω(err).Should(BeNil())
				})
			})
		*/
	})

	Describe("Simple global functions", func() {

		var (
			sm drmaa2interface.SessionManager
		)

		BeforeEach(func() {
			sm = createSessionManager()
		})

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
		var (
			sm drmaa2interface.SessionManager
		)

		BeforeEach(func() {
			sm = createSessionManager()
		})

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

})
