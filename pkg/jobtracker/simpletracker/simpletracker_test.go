package simpletracker_test

import (
	"github.com/dgruber/drmaa2interface"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"time"
)

var _ = Describe("Simpletracker", func() {

	Context("Basic tracker add job operations", func() {
		var tracker *JobTracker
		var t drmaa2interface.JobTemplate

		BeforeEach(func() {
			tracker = New("testsession")
			Ω(tracker).NotTo(BeNil())
			t = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
		})

		It("must be possible to create and destroy a tracker", func() {
			err := tracker.Destroy()
			Ω(err).Should(BeNil())
		})

		It("must be possible to add a job", func() {
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).Should(Equal("1"))

			jobs, errList := tracker.ListJobs()
			Ω(errList).Should(BeNil())
			Ω(len(jobs)).Should(Equal(1))

			err = tracker.Destroy()
			Ω(err).Should(BeNil())
		})

		It("must be possible to add an job array", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, errList := tracker.ListJobs()
			Ω(errList).To(BeNil())
			Ω(len(jobs)).To(Equal(10))

			err = tracker.Destroy()
			Ω(err).To(BeNil())
		})

		It("must be possible to get all job ids from an array job", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).To(BeNil())
			Ω(len(jobids)).To(Equal(10))

			Ω(jobids).Should(ContainElement(fmt.Sprintf("%s.1", jobid)))
		})

		It("must be possible to get a job info for a running job", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, errList := tracker.ListJobs()
			Ω(errList).To(BeNil())
			Ω(len(jobs)).To(Equal(10))

			err = tracker.Destroy()
			Ω(err).To(BeNil())
		})

		It("should be possible to wait until a job array task is finished", func() {
			t.Args = []string{"0.2"}
			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, errList := tracker.ListJobs()
			Ω(errList).To(BeNil())
			Ω(len(jobs)).To(Equal(10))

			err = tracker.Wait(jobs[2], time.Second*5, drmaa2interface.Done, drmaa2interface.Failed)
			Ω(err).To(BeNil())

			err = tracker.Wait(jobs[3], 0.0, drmaa2interface.Done, drmaa2interface.Failed)
			Ω(err).To(BeNil())

			err = tracker.Destroy()
			Ω(err).To(BeNil())
		})

		It("should be possible to list all job categories", func() {
			list, err := tracker.ListJobCategories()
			Ω(list).ShouldNot(BeNil())
			Ω(err).Should(BeNil())
		})

		Context("JobControl error cases", func() {
			var tracker *JobTracker
			var t drmaa2interface.JobTemplate

			BeforeEach(func() {
				tracker = New("testsession")
				Ω(tracker).NotTo(BeNil())
				t = drmaa2interface.JobTemplate{
					RemoteCommand: "/bin/sleep",
					Args:          []string{"0"},
				}
			})

			It("should error with wrong job id", func() {
				err := tracker.JobControl("123454321", "hold")
				Ω(err).ShouldNot(BeNil())
			})

			It("hold and release are not supported", func() {
				jobid, err := tracker.AddJob(t)
				Ω(err).Should(BeNil())
				err = tracker.JobControl(jobid, "hold")
				Ω(err).ShouldNot(BeNil())
				err = tracker.JobControl(jobid, "release")
				Ω(err).ShouldNot(BeNil())
			})

			It("should error with an undefined state", func() {
				jobid, err := tracker.AddJob(t)
				Ω(err).Should(BeNil())
				err = tracker.JobControl(jobid, "wrong")
				Ω(err).ShouldNot(BeNil())
			})
		})

	})

	Context("Basic tracker job info and modification operations", func() {
		var tracker *JobTracker
		var t drmaa2interface.JobTemplate

		BeforeEach(func() {
			tracker = New("testsession")
			Ω(tracker).NotTo(BeNil())
			t = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
		})

		It("must be in running and in done state", func() {
			t.Args = []string{"0.1"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))
			tracker.Wait(jobid, 0.0, drmaa2interface.Done)
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))
		})

		It("must be possible to start, suspend, resume, and kill a job", func() {
			t.Args = []string{"1234"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))

			err = tracker.JobControl(jobid, "suspend")
			Ω(err).Should(BeNil())

			Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Suspended))

			err = tracker.JobControl(jobid, "resume")
			Ω(err).Should(BeNil())

			Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))

			err = tracker.JobControl(jobid, "terminate")
			Ω(err).Should(BeNil())

			tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)

			Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Failed))
		})

		It("must be possible to AddJob() and DeleteJob()", func() {
			t.Args = []string{"1234"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			_, err = tracker.JobInfo(jobid)
			Ω(err).Should(BeNil())

			err = tracker.DeleteJob(jobid)
			Ω(err).ShouldNot(BeNil())

			err = tracker.JobControl(jobid, "terminate")
			Ω(err).Should(BeNil())

			tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)

			err = tracker.DeleteJob(jobid)
			Ω(err).Should(BeNil())

			err = tracker.DeleteJob(jobid)
			Ω(err).ShouldNot(BeNil())

			_, err = tracker.JobInfo(jobid)
			Ω(err).ShouldNot(BeNil())
		})

		It("must return an undetermined state for a non-existing job", func() {
			t.Args = []string{"0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, time.Second*10, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			err = tracker.DeleteJob(jobid)
			Ω(err).Should(BeNil())

			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Undetermined))
			Ω(tracker.JobState("1231231201")).Should(Equal(drmaa2interface.Undetermined))
		})

	})

	Context("Wait() operations", func() {
		var tracker *JobTracker
		var t drmaa2interface.JobTemplate

		BeforeEach(func() {
			tracker = New("testsession")
			Ω(tracker).NotTo(BeNil())
			t = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
		})

		It("must be possible to wait for a job end state", func() {
			t.Args = []string{"0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, time.Second*10, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))
		})

		It("must error when job is not found", func() {
			err := tracker.Wait("12344321", time.Second*1, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).ShouldNot(BeNil())
		})

		It("must be possible to wait for a job end state infinitely", func() {
			t.Args = []string{"0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))
		})

		It("must be possible to wait for a job when it is finished already", func() {
			t.Args = []string{"0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))

			// wait for end state when end state is already reached
			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))

			// wait for an non end state
			err = tracker.Wait(jobid, 0.0, drmaa2interface.Running, drmaa2interface.Suspended)
			Ω(err).ShouldNot(BeNil())

			// no state given
			err = tracker.Wait(jobid, 0.0)
			Ω(err).ShouldNot(BeNil())
		})

		It("must error when state is not reachable", func() {
			// wait for end state when end state is already reached
			t.Args = []string{"5"}
			jobid, err := tracker.AddJob(t)
			err = tracker.Wait(jobid, time.Millisecond*400, drmaa2interface.Suspended)
			Ω(err).ShouldNot(BeNil())
			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))
		})

		It("should be possible to wait for all job states", func() {
			// wait for end state when end state is already reached
			t.Args = []string{"0.0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			err = tracker.Wait(jobid, time.Millisecond*500, drmaa2interface.Running)
			Ω(err).Should(BeNil())
			err = tracker.Wait(jobid, time.Millisecond*500, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})
	})

	Context("JobInfo related", func() {
		var tracker *JobTracker
		var t drmaa2interface.JobTemplate

		BeforeEach(func() {
			tracker = New("testsession")
			Ω(tracker).NotTo(BeNil())
			t = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
		})

		It("should return the JobInfo object for a finished job", func() {
			t.Args = []string{"0"}
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			ji, err := tracker.JobInfo(jobid)
			Ω(err).Should(BeNil())
			Ω(ji.ID).Should(Equal(jobid))
		})

		It("should return the JobInfo for a queued job array job", func() {
			t.Args = []string{"0.2"}
			jobid, _ := tracker.AddArrayJob(t, 1, 5, 1, 1)
			id := fmt.Sprintf("%s.5", jobid)
			Ω(tracker.JobState(id)).Should(Equal(drmaa2interface.Queued))
			info, err := tracker.JobInfo(id)
			Ω(err).Should(BeNil())
			Ω(info.State).Should(Equal(drmaa2interface.Queued))
			Ω(info.Slots).Should(BeNumerically("==", 1))
			Ω(info.ID).Should(Equal(id))
			tracker.Wait(id, 0.0, drmaa2interface.Done)
		})
	})

	Context("Basic error cases", func() {
		It("must fail to add a job when JobTemplate is not correct", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{})
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

		It("must fail to add an array job when JobTemplate is not correct", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobid, err := tracker.AddArrayJob(drmaa2interface.JobTemplate{}, 1, 9, 2, 0)
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

		It("must fail to list jobs of a job array when job ID is wrong", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobids, err := tracker.ListArrayJobs("77")
			Ω(err).ShouldNot(BeNil())
			Ω(jobids).Should(BeNil())
		})

		It("must fail to list jobs of a job array when job ID is not a job array", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{RemoteCommand: "/bin/sleep"})
			Ω(err).Should(BeNil())
			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).ShouldNot(BeNil())
			Ω(jobids).Should(BeNil())
		})

		It("must fail to list jobs of a job array when job ID is wrong", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobids, err := tracker.ListArrayJobs("77")
			Ω(err).ShouldNot(BeNil())
			Ω(jobids).Should(BeNil())
		})

		It("should fail to wait for a finished job when it is in a different end state", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"0"}})
			Ω(err).Should(BeNil())
			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			state := tracker.JobState(jobid)
			Ω(state).Should(Equal(drmaa2interface.Done))
			err = tracker.Wait(jobid, 0.0, drmaa2interface.Failed)
			Ω(err).ShouldNot(BeNil())
		})

	})

	Context("Job array concurrency", func() {
		var tracker *JobTracker
		var t drmaa2interface.JobTemplate

		BeforeEach(func() {
			tracker = New("testsession")
			Ω(tracker).NotTo(BeNil())
			t = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"0.2"},
			}
		})

		runningAt := func(t time.Time, dt, ft []time.Time) int {
			running := 0
			for i := range dt {
				if t.Before(dt[i]) || t.After(ft[i]) {
					continue
				}
				running++
			}
			return running
		}

		maxParallel := func(jobids []string) int {
			max := 0
			dt := make([]time.Time, len(jobids), len(jobids))
			ft := make([]time.Time, len(jobids), len(jobids))
			for i := 0; i < len(jobids); i++ {
				info, err := tracker.JobInfo(jobids[i])
				Ω(err).Should(BeNil())
				dt[i] = info.DispatchTime
				ft[i] = info.FinishTime
			}
			for i := 0; i < len(jobids); i++ {
				if running := runningAt(dt[i], dt, ft); running > max {
					max = running
				}
				if running := runningAt(ft[i], dt, ft); running > max {
					max = running
				}
			}
			return max
		}

		It("should run all jobs sequentially if max parallel is 1", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 9, 1, 1)
			Ω(err).Should(BeNil())
			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).Should(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 9))
			tracker.Wait(jobids[8], 0.0, drmaa2interface.Done)
			Ω(maxParallel(jobids)).Should(BeNumerically("==", 1))
		})

		It("should run bunches of 3 jobs when max parallel is 3", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 9, 1, 3)
			Ω(err).Should(BeNil())
			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).Should(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 9))
			tracker.Wait(jobids[8], 0.0, drmaa2interface.Done)
			Ω(maxParallel(jobids)).Should(BeNumerically("==", 3))
		})

		It("should run all jobs parallel when max parallel is 0 or amount of jobs", func() {
			jobid, err := tracker.AddArrayJob(t, 1, 9, 1, 0)
			Ω(err).Should(BeNil())
			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).Should(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 9))
			tracker.Wait(jobids[8], 0.0, drmaa2interface.Done)
			Ω(maxParallel(jobids)).Should(BeNumerically("==", 9))

			jobid, err = tracker.AddArrayJob(t, 1, 9, 1, 9)
			Ω(err).Should(BeNil())
			jobids, err = tracker.ListArrayJobs(jobid)
			Ω(err).Should(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 9))
			tracker.Wait(jobids[8], 0.0, drmaa2interface.Done)
			Ω(maxParallel(jobids)).Should(BeNumerically("==", 9))

			jobid, err = tracker.AddArrayJob(t, 1, 9, 0, 9)
			Ω(err).Should(BeNil())
			jobids, err = tracker.ListArrayJobs(jobid)
			Ω(err).Should(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 9))
			tracker.Wait(jobids[8], 0.0, drmaa2interface.Done)
			Ω(maxParallel(jobids)).Should(BeNumerically("==", 9))
		})

		It("should should terminate job array tasks which are queued (blocked by maxParallel)", func() {
			t.Args = []string{"0.2"}
			jobid, err := tracker.AddArrayJob(t, 1, 9, 1, 4)
			Ω(err).Should(BeNil())
			// queued
			err = tracker.JobControl(fmt.Sprintf("%s.8", jobid), "terminate")
			Ω(err).Should(BeNil())
			// running
			err = tracker.JobControl(fmt.Sprintf("%s.1", jobid), "terminate")
			Ω(err).Should(BeNil())
			err = tracker.Wait(fmt.Sprintf("%s.8", jobid), 0.0, drmaa2interface.Failed)
			Ω(err).Should(BeNil())
			err = tracker.Wait(fmt.Sprintf("%s.9", jobid), 0.0, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("must fail to suspend/resume a job array task which is queued", func() {
			t.Args = []string{"0.1"}
			jobid, err := tracker.AddArrayJob(t, 1, 4, 1, 1)
			Ω(err).Should(BeNil())
			err = tracker.JobControl(fmt.Sprintf("%s.4", jobid), "suspend")
			Ω(err).ShouldNot(BeNil())
			err = tracker.JobControl(fmt.Sprintf("%s.4", jobid), "resume")
			Ω(err).ShouldNot(BeNil())
		})

	})

})
