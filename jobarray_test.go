package drmaa2os

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/fakes"
)

var _ = Describe("Jobarray", func() {

	Context("Basic JobArray Operations", func() {

		id := "ID"
		session := "sessionname"
		template := drmaa2interface.JobTemplate{
			RemoteCommand: "Command",
			JobName:       "name",
		}

		It("should return all set content", func() {
			aj := newArrayJob(id, session, template, nil)
			Ω(aj).ShouldNot(BeNil())
			Ω(aj.GetID()).Should(Equal(id))
			Ω(aj.GetSessionName()).Should(Equal(session))
			Ω(aj.GetJobTemplate()).ShouldNot(BeNil())
			Ω(aj.GetJobs()).Should(BeNil())
		})

		It("should perform operation on all jobs", func() {
			js := []drmaa2interface.Job{
				&fakes.Job{
					ID:      "1",
					Session: "session",
				},
				&fakes.Job{
					ID:      "2",
					Session: "session",
				}}

			aj := newArrayJob(id, session, template, js)
			Ω(aj).ShouldNot(BeNil())

			aj.Suspend()
			jobs := aj.GetJobs()
			Ω(jobs).ShouldNot(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Suspended))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Suspended))

			aj.Resume()
			jobs = aj.GetJobs()
			Ω(jobs).ShouldNot(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Running))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Running))

			aj.Hold()
			jobs = aj.GetJobs()
			Ω(jobs).ShouldNot(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.QueuedHeld))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.QueuedHeld))

			aj.Release()
			jobs = aj.GetJobs()
			Ω(jobs).ShouldNot(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Running))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Running))

			aj.Terminate()
			jobs = aj.GetJobs()
			Ω(jobs).ShouldNot(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Failed))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Failed))
		})

	})

})
