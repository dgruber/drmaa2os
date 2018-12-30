package drmaa2os

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletrackerfakes"
	"time"
)

const tempdb string = "drmaa2ostest.db"

var _ = Describe("Jobsession", func() {

	Describe("standard operations", func() {
		var (
			js drmaa2interface.JobSession
			jt drmaa2interface.JobTemplate
		)

		BeforeEach(func() {
			js = NewJobSession("testsession", []jobtracker.JobTracker{simpletracker.New("testsession")})
			jt = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"0.1"},
			}
		})

		It("should return the job session name", func() {
			name, err := js.GetSessionName()
			Ω(err).Should(BeNil())
			Ω(name).Should(Equal("testsession"))
			Ω(js.Close()).Should(BeNil())
		})

		It("should be to get the contact string", func() {
			_, err := js.GetContact()
			Ω(err).Should(BeNil())
		})

		It("should be to get the job categories", func() {
			js = NewJobSession("testsession", []jobtracker.JobTracker{simpletrackerfakes.New("testsession")})
			categories, err := js.GetJobCategories()
			Ω(err).Should(BeNil())
			Ω(categories).ShouldNot(BeNil())
			Ω(len(categories)).Should(BeNumerically("==", 2))
		})

		It("should be able to submit a job and get access to it", func() {
			job, err := js.RunJob(jt)
			Ω(err).Should(BeNil())

			template, errTempl := job.GetJobTemplate()
			Ω(errTempl).Should(BeNil())
			Ω(template).Should(Equal(jt))

			jobs, errJobs := js.GetJobs(drmaa2interface.CreateJobInfo())
			Ω(errJobs).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 1))
		})

		It("should be able to wait for a started job", func() {
			job, err := js.RunJob(jt)
			Ω(err).Should(BeNil())

			jobid := job.GetID()

			var jobs []drmaa2interface.Job
			jobs = append(jobs, job)

			j, err := js.WaitAnyStarted(jobs, time.Second*2)
			Ω(err).Should(BeNil())
			Ω(j.GetID()).Should(Equal(jobid))
			//Ω(j.GetState()).Should(Equal(drmaa2interface.Failed))
			Ω(js.Close()).Should(BeNil())
		})

		It("should be able to wait for a finished job", func() {
			job, err := js.RunJob(jt)
			Ω(err).Should(BeNil())

			jobid := job.GetID()

			j, err := js.WaitAnyTerminated([]drmaa2interface.Job{job}, time.Second*2)
			Ω(err).Should(BeNil())
			Ω(j.GetID()).Should(Equal(jobid))
			//Ω(j.GetState()).Should(Equal(drmaa2interface.Done))
			Ω(js.Close()).Should(BeNil())
		})
	})

	Describe("waitAny with fakes", func() {

		It("should return when one job is running", func() {
			job1 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "", time.Millisecond*100)
			job2 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "", time.Millisecond*5000)

			var array []drmaa2interface.Job
			array = append(array, job1)
			array = append(array, job2)

			job, err := waitAny(true, array, time.Second*4)

			Ω(err).Should(BeNil())
			Ω(job.GetState()).Should(Equal(drmaa2interface.Running))
		})

		It("should return with an error when timeout is reached", func() {
			job1 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "", time.Millisecond*2000)
			job2 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "", time.Millisecond*1500)

			var array []drmaa2interface.Job
			array = append(array, job1)
			array = append(array, job2)

			job, err := waitAny(true, array, time.Second*1)

			Ω(err).ShouldNot(BeNil())
			Ω(job).Should(BeNil())
		})

		It("should return with an error when all job wait calls errors immediately", func() {
			job1 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "immediate error", time.Millisecond*2000)
			job2 := simpletrackerfakes.NewFakeJob(drmaa2interface.Running, "immediate error", time.Millisecond*1500)

			var array []drmaa2interface.Job
			array = append(array, job1)
			array = append(array, job2)

			job, err := waitAny(true, array, time.Second*1)

			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(errors.New("Error waiting for jobs")))
			Ω(job).Should(BeNil())
		})
	})

	Describe("basic job array functionality", func() {
		var (
			js drmaa2interface.JobSession
			jt drmaa2interface.JobTemplate
		)

		BeforeEach(func() {
			js = NewJobSession("testsession", []jobtracker.JobTracker{simpletracker.New("testsession")})
			jt = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"0.1"},
			}
		})

		It("should be possible to submit a job array (bulk job)", func() {
			arrayjob, err := js.RunBulkJobs(jt, 1, 10, 1, 2)
			Ω(err).Should(BeNil())

			jobid := arrayjob.GetID()
			Ω(jobid).ShouldNot(Equal(""))

			jobs := arrayjob.GetJobs()
			Ω(len(jobs)).Should(Equal(10))

			j, err := js.WaitAnyTerminated(jobs, time.Second*20)
			Ω(err).Should(BeNil())
			Ω(j.GetID()).Should(ContainSubstring(jobid))
			Ω(j.GetState()).Should(Equal(drmaa2interface.Done))
			Ω(js.Close()).Should(BeNil())
		})

		It("should be possible to terminate a job array (bulk job)", func() {
			jt.Args = []string{"100"}

			arrayjob, err := js.RunBulkJobs(jt, 1, 10, 1, 5)
			Ω(err).Should(BeNil())

			jobid := arrayjob.GetID()
			Ω(jobid).ShouldNot(Equal(""))

			err = arrayjob.Terminate()
			Ω(err).Should(BeNil())

			for _, j := range arrayjob.GetJobs() {
				Ω(j.GetState()).Should(Equal(drmaa2interface.Failed))
			}
			Ω(js.Close()).Should(BeNil())
		})

	})

})
