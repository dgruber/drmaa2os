package drmaa2os

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/fakes"
)

var _ = Describe("JobarrayHlp", func() {

	Context("Perform actions without issues", func() {

		var jobs []drmaa2interface.Job

		BeforeEach(func() {
			jobs = []drmaa2interface.Job{
				&fakes.Job{
					ID:      "1",
					Session: "session",
				},
				&fakes.Job{
					ID:      "2",
					Session: "session",
				}}
		})

		It("should suspend all jobs", func() {
			err := jobAction(suspend, jobs)
			Ω(err).Should(BeNil())
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Suspended))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Suspended))
		})

		It("should resume all jobs", func() {
			err := jobAction(resume, jobs)
			Ω(err).Should(BeNil())
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Running))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Running))
		})

		It("should hold all jobs", func() {
			err := jobAction(hold, jobs)
			Ω(err).Should(BeNil())
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.QueuedHeld))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.QueuedHeld))
		})

		It("should release all jobs", func() {
			err := jobAction(release, jobs)
			Ω(err).Should(BeNil())
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Running))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Running))
		})

		It("should terminate all jobs", func() {
			err := jobAction(terminate, jobs)
			Ω(err).Should(BeNil())
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Failed))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Failed))
		})

	})

	Context("Perform action with errors", func() {

		var jobs []drmaa2interface.Job

		BeforeEach(func() {
			jobs = append(jobs, &fakes.Job{
				ID:               "1",
				Session:          "session",
				ErrorWhenSuspend: true,
			})
			jobs = append(jobs, &fakes.Job{
				ID:      "2",
				Session: "session",
			})
		})

		It("should report an error when first job fails to get suspended", func() {
			err := jobAction(suspend, jobs)
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(Equal("Job 1 error: Some error happened "))
			Ω(jobs[1].GetState()).Should(Equal(drmaa2interface.Suspended))
		})

		It("should report an error when second job fails to get suspended", func() {
			jobs = []drmaa2interface.Job{
				&fakes.Job{
					ID:               "1",
					Session:          "session",
					ErrorWhenSuspend: false,
				},
				&fakes.Job{
					ID:               "2",
					Session:          "session",
					ErrorWhenSuspend: true,
				},
			}

			err := jobAction(suspend, jobs)
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(Equal("Job 2 error: Some error happened "))
			Ω(jobs[0].GetState()).Should(Equal(drmaa2interface.Suspended))
		})

	})

})
