package singularity_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/singularity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Singularity", func() {

	var template drmaa2interface.JobTemplate

	BeforeEach(func() {
		template = drmaa2interface.JobTemplate{
			RemoteCommand: "/bin/sleep",
			Args:          []string{"1"},
			JobCategory:   "test.simg",
			OutputPath:    "/dev/stdout",
			ErrorPath:     "/dev/stderr",
		}
	})

	Context("Happy Path", func() {

		It("should create a Singularity session", func() {
			_, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
		})

		It("should create a new Singularity container", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))
			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("should create a new Singularity container with Singularity arguments specified in JobTemplate", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			template.ExtensionList = map[string]string{
				"pid": "",
			}
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))
			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("should list the Singularity containers", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))
			job2, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job2).ShouldNot(Equal(""))

			jobs, err := st.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			err = st.Wait(job2, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("should create a an array of Singularity containers", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			ajID, err := st.AddArrayJob(template, 1, 2, 1, 1)
			Ω(err).Should(BeNil())
			Ω(ajID).ShouldNot(Equal(""))

			jobs, err := st.ListArrayJobs(ajID)
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 2))

			err = st.Wait(jobs[0], drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			err = st.Wait(jobs[1], drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("should be able to suspend and resume a Singularity container", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			template.RemoteCommand = "/bin/sleep"
			template.Args = []string{"1000"}
			template.ExtensionList = map[string]string{
				"pid": "",
			}
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))

			st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Running)

			err = st.JobControl(job, "suspend")
			Ω(err).Should(BeNil())
			Ω(st.JobState(job)).Should(BeNumerically("==", drmaa2interface.Suspended))

			err = st.JobControl(job, "resume")
			Ω(err).Should(BeNil())
			Ω(st.JobState(job)).Should(BeNumerically("==", drmaa2interface.Running))

			err = st.JobControl(job, "terminate")
			Ω(err).Should(BeNil())

			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
		})

		It("should create a JobInfo describing the Singularity container", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))

			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			jobInfo, err := st.JobInfo(job)
			Ω(err).Should(BeNil())
			Ω(jobInfo.ID).Should(Equal(job))
			Ω(jobInfo.State).Should(Equal(drmaa2interface.Done))
		})

		It("should list the JobCategories which are the container images which is unknown and 0", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			jcats, err := st.ListJobCategories()
			Ω(err).Should(BeNil())
			Ω(len(jcats)).Should(BeNumerically("==", 0))
		})

		It("should delete the job when the job is an end state", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			job, err := st.AddJob(template)
			Ω(err).Should(BeNil())
			Ω(job).ShouldNot(Equal(""))
			err = st.DeleteJob(job)
			// job is running therfore reaping should not be possible
			Ω(err).ShouldNot(BeNil())
			err = st.Wait(job, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			// finished jobs should be able to be reaped from DB
			err = st.DeleteJob(job)
			Ω(err).Should(BeNil())
		})

	})

	Context("Basic error cases", func() {

		It("should fail to create a new Singularity container when image is missing", func() {
			st, err := New("singularity_test_session")
			Ω(err).Should(BeNil())
			t2 := template
			t2.JobCategory = ""
			jobid, err := st.AddJob(t2)
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

	})

})
