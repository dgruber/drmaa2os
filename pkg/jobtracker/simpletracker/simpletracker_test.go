package simpletracker_test

import (
	"github.com/dgruber/drmaa2interface"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("Simpletracker", func() {

	Context("Basic tracker add job operations", func() {

		It("must be possible to create and destroy a tracker", func() {
			tracker := New("testsession")
			Ω(tracker).ShouldNot(BeNil())
			err := tracker.Destroy()
			Ω(err).Should(BeNil())
		})

		It("must be possible to add a job", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())

			t := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"2"},
			}

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
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())

			t := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"2"},
			}

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
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())

			t := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"2"},
			}

			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobids, err := tracker.ListArrayJobs(jobid)
			Ω(err).To(BeNil())
			Ω(len(jobids)).To(Equal(10))

			Ω(jobids).Should(ContainElement(fmt.Sprintf("%s.1", jobid)))
		})

		It("must be possible to get a job info for a running job", func() {
			tracker := New("testsession")
			Ω(tracker).NotTo(BeNil())

			t := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"2"},
			}

			jobid, err := tracker.AddArrayJob(t, 1, 10, 1, 0)
			Ω(err).To(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, errList := tracker.ListJobs()
			Ω(errList).To(BeNil())
			Ω(len(jobs)).To(Equal(10))

			err = tracker.Destroy()
			Ω(err).To(BeNil())

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
			jobid, err := tracker.AddJob(t)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			Ω(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))

			// TODO must be done after a second
			// Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Done))
		})

		It("must be possible to start, suspend, resume, and kill a job", func() {
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

			Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Failed))

			// TODO must be done after a second
			// Eventually(tracker.JobState(jobid)).Should(Equal(drmaa2interface.Running))

		})

	})

})
