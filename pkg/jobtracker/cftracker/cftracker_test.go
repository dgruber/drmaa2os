package cftracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Cftracker", func() {

	Context("Basic Operations", func() {
		var client *cftracker
		jt := drmaa2interface.JobTemplate{
			JobName:       "name",
			RemoteCommand: "command",
			Args:          []string{"123"},
			JobCategory:   "guid",
		}

		BeforeEach(func() {
			client = newFake("addr", "username", "password", "jobsession")
			Ω(client).ShouldNot(BeNil())
		})

		It("should be possible to list jobs", func() {
			jobs, err := client.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 1))
			Ω(jobs[0]).Should(Equal("GUID"))
		})

		It("should be possible to add a job", func() {
			jobid, err := client.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).Should(Equal("GUID"))
		})

		It("should be possible to add an array job", func() {
			arrayjobid, err := client.AddArrayJob(jt, 1, 10, 1, 0)
			Ω(err).Should(BeNil())
			ids, err := client.ListArrayJobs(arrayjobid)
			Ω(err).Should(BeNil())
			Ω(len(ids)).Should(Equal(10))
			for i := 0; i < 10; i++ {
				Ω(ids[i]).Should(Equal("GUID"))
			}
		})

		It("should show the job state", func() {
			state := client.JobState("GUID")
			Ω(state).Should(Equal(drmaa2interface.Failed))

			state = client.JobState("PENDING")
			Ω(state).Should(Equal(drmaa2interface.Queued))

			state = client.JobState("RUNNING")
			Ω(state).Should(Equal(drmaa2interface.Running))

			state = client.JobState("CANCELING")
			Ω(state).Should(Equal(drmaa2interface.Running))

			state = client.JobState("SUCCEEDED")
			Ω(state).Should(Equal(drmaa2interface.Done))

			state = client.JobState("unknown")
			Ω(state).Should(Equal(drmaa2interface.Undetermined))

			state = client.JobState("error")
			Ω(state).Should(Equal(drmaa2interface.Undetermined))

		})

		It("should show the JobInfo", func() {
			ji, err := client.JobInfo("GUID")
			Ω(err).Should(BeNil())
			Ω(ji.State).Should(Equal(drmaa2interface.Failed))
		})

		It("should be possible to control the job", func() {
			// unsupported
			err := client.JobControl("GUID", "suspend")
			Ω(err).ShouldNot(BeNil())
			err = client.JobControl("GUID", "resume")
			Ω(err).ShouldNot(BeNil())
			err = client.JobControl("GUID", "hold")
			Ω(err).ShouldNot(BeNil())
			err = client.JobControl("GUID", "release")
			Ω(err).ShouldNot(BeNil())
			// supported
			err = client.JobControl("noerror", "terminate")
			Ω(err).Should(BeNil())
		})

		It("should be possible to delete the job", func() {
			// purge task info not possible
			err := client.DeleteJob("1")
			Ω(err).ShouldNot(BeNil())
		})

		It("should be possible to list job categories (apps)", func() {
			cats, err := client.ListJobCategories()
			Ω(err).Should(BeNil())
			Ω(cats).ShouldNot(BeNil())
			Ω(len(cats)).Should(BeNumerically("==", 1))
			Ω(cats[0]).Should(Equal("guid"))
		})

	})

	Context("Expected Errors", func() {
		var client *cftracker
		jt := drmaa2interface.JobTemplate{
			JobName:       "name",
			RemoteCommand: "command",
			Args:          []string{"123"},
		}

		BeforeEach(func() {
			client = newFake("addr", "username", "password", "jobsession")
			Ω(client).ShouldNot(BeNil())
		})

		It("should error when a broken job template is used in AddJob", func() {
			jt.RemoteCommand = ""
			_, err := client.AddJob(jt)
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when a broken job template is used in AddArrayJob", func() {
			jt.RemoteCommand = ""
			_, err := client.AddArrayJob(jt, 1, 10, 1, 0)
			Ω(err).ShouldNot(BeNil())
		})

		It("should error when creating a job in AddJob fails", func() {
			jt.RemoteCommand = "error"
			_, err := client.AddJob(jt)
			Ω(err).ShouldNot(BeNil())
		})

		It("JobInfo should fail when task is not found", func() {
			_, err := client.JobInfo("error")
			Ω(err).ShouldNot(BeNil())
		})

		It("AddArrayJob should fail when range is wrong", func() {
			_, err := client.AddArrayJob(jt, 1, 10, 0, 0)
			Ω(err).ShouldNot(BeNil())
		})

		It("JobControl should fail when undefined state is used", func() {
			err := client.JobControl("GUID", "unknownState")
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("Not yet implemented", func() {
		var client *cftracker

		BeforeEach(func() {
			client = newFake("addr", "username", "password", "jobsession")
			Ω(client).ShouldNot(BeNil())
		})

		It("should error when DeleteJob is called", func() {
			err := client.DeleteJob("error")
			Ω(err).ShouldNot(BeNil())
		})

	})

})
