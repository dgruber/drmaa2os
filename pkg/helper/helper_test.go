package helper_test

import (
	. "github.com/dgruber/drmaa2os/pkg/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletrackerfakes"
	"time"
)

var _ = Describe("Helper", func() {

	Context("Array Job ID convert functions", func() {

		It("should generate and resolve an array job ID into job IDs", func() {
			guids := []string{"1", "2", "3"}

			id := Guids2ArrayJobID(guids)
			guidsOut, err := ArrayJobID2GUIDs(id)

			Ω(err).Should(BeNil())
			Ω(guidsOut).Should(BeEquivalentTo(guids))
		})

	})

	Context("Create array job out with single job submissions", func() {

		jt := drmaa2interface.JobTemplate{RemoteCommand: "test"}

		It("AddArrayJobAsSingleJobs should work", func() {
			fakeTracker := simpletrackerfakes.New("testsession")
			_, err := AddArrayJobAsSingleJobs(jt, fakeTracker, 11, 110, 2)
			Ω(err).Should(BeNil())
			jobs, errJobs := fakeTracker.ListJobs()
			Ω(errJobs).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 50))
		})

	})

	Context("Error cases", func() {

		It("should return nothing when array job id is not parsable", func() {
			ajid, err := ArrayJobID2GUIDs("_")
			Ω(err).ShouldNot(BeNil())
			Ω(ajid).Should(BeEmpty())
		})

	})

	Context("State", func() {

		It("should signal when a state is matched", func() {
			state := IsInExpectedState(drmaa2interface.Done, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(state).Should(BeTrue())
			state = IsInExpectedState(drmaa2interface.Done, drmaa2interface.Done, drmaa2interface.Failed)
			Ω(state).Should(BeTrue())
			state = IsInExpectedState(drmaa2interface.Done, drmaa2interface.Done)
			Ω(state).Should(BeTrue())
		})

		It("should signal when a state is not matched", func() {
			state := IsInExpectedState(drmaa2interface.Failed, drmaa2interface.Queued, drmaa2interface.Done)
			Ω(state).Should(BeFalse())
			state = IsInExpectedState(drmaa2interface.Failed, drmaa2interface.Done)
			Ω(state).Should(BeFalse())
			state = IsInExpectedState(drmaa2interface.Failed, drmaa2interface.Requeued)
			Ω(state).Should(BeFalse())
			state = IsInExpectedState(drmaa2interface.Failed, drmaa2interface.Undetermined)
			Ω(state).Should(BeFalse())
		})

		StateAfter := func(t jobtracker.JobTracker, id string, d time.Duration, operation string) {
			go func() {
				<-time.Tick(d)
				t.JobControl(id, operation)
			}()
		}

		WaitFor := func(tracker jobtracker.JobTracker, jobid, operation string, expectedStates ...drmaa2interface.JobState) error {
			StateAfter(tracker, jobid, time.Millisecond*50, operation)
			return WaitForState(tracker, jobid, time.Second*1, expectedStates...)
		}

		It("should block until state is reached or error", func() {
			tracker := simpletrackerfakes.New("testsession")
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{
				JobName: "testjob",
			})
			Ω(err).Should(BeNil())

			Ω(WaitFor(tracker, jobid, "terminate", drmaa2interface.Failed)).Should(BeNil())
			Ω(WaitFor(tracker, jobid, "suspend", drmaa2interface.Done, drmaa2interface.Suspended)).Should(BeNil())
			Ω(WaitFor(tracker, jobid, "terminate", drmaa2interface.Done, drmaa2interface.Failed)).Should(BeNil())

			// timeout
			Ω(WaitFor(tracker, jobid, "terminate", drmaa2interface.Suspended)).ShouldNot(BeNil())
			Ω(WaitFor(tracker, jobid, "terminate", drmaa2interface.Done, drmaa2interface.Suspended, drmaa2interface.Running)).ShouldNot(BeNil())
		})

		It("should return immediately when job is already in state or no timeout is given", func() {
			tracker := simpletrackerfakes.New("testsession")
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{
				JobName: "testjob",
			})
			Ω(err).Should(BeNil())
			tracker.JobControl(jobid, "suspend")
			Ω(WaitForState(tracker, jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Suspended)).Should(BeNil())
			Ω(WaitForState(tracker, jobid, 0.0, drmaa2interface.Done, drmaa2interface.Failed)).ShouldNot(BeNil())
		})

	})

})
