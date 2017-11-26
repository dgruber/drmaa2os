package kubernetestracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"

	"time"
)

var _ = Describe("KubernetesTracker", func() {

	Context("Basic interface test", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				//JobName:       "name1",
				RemoteCommand: "command",
				JobCategory:   "golang:latest",
				Args:          []string{"0"},
			}
			var err error
			kt, err = New()
			Ω(err).Should(BeNil())
		})

		It("should be possible to AddJob()", func() {
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
		})

		It("should be possible to AddArrayJob()", func() {
			jobid, err := kt.AddArrayJob(jt, 1, 2, 1, 0)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
		})

		It("should be possible to ListJobs()", func() {
			jobids, err := kt.ListJobs()
			Ω(err).Should(BeNil())
			Ω(jobids).ShouldNot(BeNil())
		})

		It("should be possible to ListArrayJobs()", func() {
			jobids, err := kt.ListArrayJobs("123")
			Ω(err).ShouldNot(BeNil())
			Ω(jobids).Should(BeNil())
		})

	})

	Context("Unsupported interface functions", func() {
		var kt jobtracker.JobTracker

		BeforeEach(func() {
			var err error
			kt, err = New()
			Ω(err).Should(BeNil())
		})

		It("Unsupported ListJobCategories()", func() {
			_, err := kt.ListJobCategories()
			Ω(err).Should(BeNil())
		})

		It("Unsupported JobInfo()", func() {
			_, err := kt.JobInfo("jobid")
			Ω(err).Should(BeNil())
		})

		It("Unsupported Wait()", func() {
			err := kt.Wait("jobid", time.Second*0, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("Unsupported DeleteJob()", func() {
			err := kt.DeleteJob("jobid")
			Ω(err).Should(BeNil())
		})

	})

	Context("Basic Kubernetes Job Workflow", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				RemoteCommand: "command",
				JobCategory:   "golang:latest",
				Args:          []string{"1"},
			}
			var err error
			kt, err = New()
			Ω(err).Should(BeNil())
		})

		It("Should be possible to submit and delete a job", func() {
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			<-time.After(time.Millisecond * 100)
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Running))
			err = kt.JobControl(jobid, "terminate")
			Ω(err).Should(BeNil())
			<-time.After(time.Millisecond * 100)
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Failed))
		})

	})

})
