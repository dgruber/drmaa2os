package d2hlp_test

import (
	. "github.com/dgruber/drmaa2os/pkg/d2hlp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Jinfomatcher", func() {

	Context("unset jobinfo", func() {

		var (
			e      drmaa2interface.JobInfo
			filter drmaa2interface.JobInfo
		)

		BeforeEach(func() {
			e = drmaa2interface.JobInfo{}
			filter = drmaa2interface.CreateJobInfo()
		})

		It("should never filter", func() {
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.Slots = 2
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.CPUTime = int64(2)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.DispatchTime = time.Now()
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.ExitStatus = 0
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.ExitStatus = 10
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.Annotation = "blub"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.ID = "123"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
			e.SubmissionMachine = "submission"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should display if a jobinfo struct is unset", func() {
			Ω(JobInfoIsUnset(filter)).Should(BeTrue())
			Ω(JobInfoIsUnset(e)).Should(BeFalse())
		})
	})

	Context("when single entries are for element and in filter", func() {

		var (
			e      drmaa2interface.JobInfo
			filter drmaa2interface.JobInfo
		)

		BeforeEach(func() {
			e = drmaa2interface.CreateJobInfo()
			filter = drmaa2interface.CreateJobInfo()
		})

		It("should match jobid", func() {
			filter.ID = "13"
			e.ID = "13"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match exit status", func() {
			filter.ExitStatus = 13
			e.ExitStatus = 13
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match termination signal", func() {
			filter.TerminatingSignal = "SIGKILL"
			e.TerminatingSignal = "SIGKILL"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match annotation", func() {
			filter.Annotation = "anno"
			e.Annotation = "anno"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match state", func() {
			filter.State = drmaa2interface.Failed
			e.State = drmaa2interface.Failed
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match substate", func() {
			filter.SubState = "substate"
			e.SubState = "substate"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		// GFD 231: "the job is executed on a superset of the given list of machines"

		It("should match allocated machines", func() {
			filter.AllocatedMachines = []string{"ubuntu"}
			e.AllocatedMachines = []string{"ubuntu"}
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match superset of allocated machines", func() {
			filter.AllocatedMachines = []string{"ubuntu", "suse"}
			e.AllocatedMachines = []string{"ubuntu"}
			Ω(JobInfoMatches(e, filter)).ShouldNot(BeTrue())
		})

		It("should not match subset of allocated machines", func() {
			filter.AllocatedMachines = []string{"ubuntu"}
			e.AllocatedMachines = []string{"ubuntu", "suse"}
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should not match other set of allocated machines", func() {
			filter.AllocatedMachines = []string{"ubuntu2", "suse"}
			e.AllocatedMachines = []string{"ubuntu"}
			Ω(JobInfoMatches(e, filter)).ShouldNot(BeTrue())
		})

		It("should match submission machine", func() {
			filter.SubmissionMachine = "ubuntu"
			e.SubmissionMachine = "ubuntu"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match job owner", func() {
			filter.JobOwner = "root"
			e.JobOwner = "root"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match slots", func() {
			filter.Slots = 12
			e.Slots = 12
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match queue name", func() {
			filter.QueueName = "queuename"
			e.QueueName = "queuename"
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match wallclock time", func() {
			filter.WallclockTime = time.Minute * 1
			e.WallclockTime = time.Minute * 1
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match wallclock time (at least)", func() {
			filter.WallclockTime = time.Minute * 1
			e.WallclockTime = time.Minute * 2
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match cpu time", func() {
			filter.CPUTime = 60
			e.CPUTime = 60
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match cpu time (at least)", func() {
			filter.CPUTime = 60
			e.CPUTime = 120
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match submission time", func() {
			filter.SubmissionTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			e.SubmissionTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match submission time when job was submitted after given time", func() {
			filter.SubmissionTime = time.Date(2015, 12, 31, 23, 59, 59, 0, time.Local)
			e.SubmissionTime = time.Date(2016, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match dispatch time", func() {
			filter.DispatchTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			e.DispatchTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match dispatch time when the job was dispatched after given time", func() {
			filter.DispatchTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			e.DispatchTime = time.Date(2001, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

		It("should match finish time", func() {
			filter.FinishTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			e.FinishTime = time.Date(1999, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeTrue())
		})

	})

	Context("when single entries are for element and in filter but they don't match", func() {

		var (
			e      drmaa2interface.JobInfo
			filter drmaa2interface.JobInfo
		)

		BeforeEach(func() {
			e = drmaa2interface.CreateJobInfo()
			filter = drmaa2interface.CreateJobInfo()
		})

		It("should not match jobid", func() {
			filter.ID = "13"
			e.ID = "0"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match exit status", func() {
			filter.ExitStatus = 13
			e.ExitStatus = 14
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match termination signal", func() {
			filter.TerminatingSignal = "SIGKILL"
			e.TerminatingSignal = "SIGCONT"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match annotation", func() {
			filter.Annotation = "anno"
			e.Annotation = "annodazumal"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match state", func() {
			filter.State = drmaa2interface.Failed
			e.State = drmaa2interface.Done
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match substate", func() {
			filter.SubState = "substate"
			e.SubState = "othersubstate"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match allocated machines", func() {
			filter.AllocatedMachines = []string{"ubuntu"}
			e.AllocatedMachines = []string{"suse"}
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match submission machine", func() {
			filter.SubmissionMachine = "ubuntu"
			e.SubmissionMachine = "suse"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match job owner", func() {
			filter.JobOwner = "root"
			e.JobOwner = "user"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match slots", func() {
			filter.Slots = 12
			e.Slots = 13
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match queue name", func() {
			filter.QueueName = "queuename"
			e.QueueName = "noqueue"
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match wallclock time", func() {
			filter.WallclockTime = time.Minute * 1
			e.WallclockTime = time.Second * 1
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match cpu time", func() {
			filter.CPUTime = 60
			e.CPUTime = 30
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match submission time", func() {
			filter.SubmissionTime = time.Date(2001, 12, 31, 23, 59, 59, 0, time.Local)
			e.SubmissionTime = time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match dispatch time", func() {
			filter.DispatchTime = time.Date(2002, 12, 31, 23, 59, 59, 0, time.Local)
			e.DispatchTime = time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

		It("should not match finish time", func() {
			filter.FinishTime = time.Date(2002, 12, 31, 23, 59, 59, 0, time.Local)
			e.FinishTime = time.Date(2000, 12, 31, 23, 59, 59, 0, time.Local)
			Ω(JobInfoMatches(e, filter)).Should(BeFalse())
		})

	})

	Context("StringFilter", func() {

		It("should filter strings", func() {
			filter := []string{"1", "3", "5", "7"}
			f := NewStringFilter(filter)
			Ω(f.IsIncluded("1")).Should(BeTrue())
			Ω(f.IsIncluded("2")).Should(BeFalse())
			Ω(f.IsIncluded("3")).Should(BeTrue())
			Ω(f.IsIncluded("4")).Should(BeFalse())
			Ω(f.IsIncluded("5")).Should(BeTrue())
			Ω(f.IsIncluded("6")).Should(BeFalse())
		})

		It("should not filter strings", func() {
			f := NewStringFilter(nil)
			Ω(f.IsIncluded("1")).Should(BeFalse())
			Ω(f.IsIncluded("2")).Should(BeFalse())
			Ω(f.IsIncluded("3")).Should(BeFalse())
		})

	})

})
