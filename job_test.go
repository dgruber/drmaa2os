package drmaa2os

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/dgruber/drmaa2interface"
	_ "github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletrackerfakes"
)

var _ = Describe("Job", func() {

	Context("creation and destruction", func() {
		var job drmaa2interface.Job
		var template drmaa2interface.JobTemplate
		var tracker *simpletrackerfakes.JobTracker

		BeforeEach(func() {
			template = drmaa2interface.JobTemplate{JobName: "jobname"}
			tracker = simpletrackerfakes.New("jobsession")
			tracker.AddJob(template)
			job = newJob("13", "jobsession", template, tracker)
		})

		It("should return the job specific details", func() {
			Ω(job.GetID()).Should(Equal("13"))
			Ω(job.GetSessionName()).Should(Equal("jobsession"))

			jtemplate, err := job.GetJobTemplate()
			Ω(err).Should(BeNil())
			Ω(jtemplate).Should(Equal(template))

			_, jiError := job.GetJobInfo()
			Ω(jiError).ShouldNot(BeNil())
		})

		It("should be in the expected state", func() {
			err := job.Hold()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.QueuedHeld))

			err = job.Release()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Running))

			err = job.Suspend()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Suspended))

			err = job.Resume()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Running))

			err = job.Release()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Running))

			err = job.Terminate()
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Failed))
		})

		It("should be able to wait for a state", func() {
			err := job.WaitStarted(time.Second * 1)
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Running))

			err = job.WaitTerminated(time.Second * 1)
			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Done))
		})

		It("should reap the job", func() {
			err := job.Reap()
			// should error since not in an end state
			Ω(err).ShouldNot(BeNil())

			err = job.Terminate()
			Ω(err).Should(BeNil())

			err = job.Reap()
			Ω(err).Should(BeNil())
		})
	})

})
