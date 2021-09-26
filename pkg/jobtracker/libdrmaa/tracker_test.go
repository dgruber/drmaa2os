package libdrmaa

import (
	"time"

	"github.com/dgruber/drmaa2interface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tracker", func() {

	var sleeperJob drmaa2interface.JobTemplate

	BeforeEach(func() {
		sleeperJob = drmaa2interface.JobTemplate{
			RemoteCommand: "sleep",
			Args:          []string{"0"},
			JobName:       "mandatory",
		}
	})

	Context("Basic workflow", func() {

		It("should run a single job", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			jobID, err := d.AddJob(sleeperJob)
			Expect(err).To(BeNil())
			Expect(jobID).NotTo(Equal(""))

			err = d.Wait(jobID, time.Second*31, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			state, substate, err := d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(substate).To(Equal(""))
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))

			jobs, err := d.ListJobs()
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 1))
			Expect(jobs[0]).To(Equal(jobID))

			state, _, err = d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(state).To(Equal(drmaa2interface.Done))

			jobInfo, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())
			Expect(jobInfo.ID).To(Equal(jobID))
			Expect(jobInfo.State).To(Equal(drmaa2interface.Done))
		})

		It("should run a job array", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			arrayJobID, err := d.AddArrayJob(sleeperJob, 1, 10, 1, 0)
			Expect(err).To(BeNil())
			Expect(arrayJobID).NotTo(Equal(""))

			ids, err := d.ListArrayJobs(arrayJobID)
			Expect(err).To(BeNil())
			Expect(len(ids)).To(BeNumerically("==", 10))
			for _, id := range ids {
				err = d.Wait(id, time.Second*31, drmaa2interface.Done, drmaa2interface.Failed)
				Expect(err).To(BeNil())

				state, substate, err := d.JobState(id)
				Expect(err).To(BeNil())
				Expect(substate).To(Equal(""))
				Expect(state.String()).To(Equal(drmaa2interface.Done.String()))
			}

			jobs, err := d.ListJobs()
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 10))

			state, _, err := d.JobState(ids[9])
			Expect(err).To(BeNil())
			Expect(state).To(Equal(drmaa2interface.Done))

			jobInfo, err := d.JobInfo(ids[9])
			Expect(err).To(BeNil())
			Expect(jobInfo.State).To(Equal(drmaa2interface.Done))
		})

	})

	Context("Failing jobs", func() {

		It("should mark a job as done and show exit status when exit status != 0", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			sleeperJob.RemoteCommand = "/bin/bash"
			sleeperJob.Args = []string{"-c", `exit 1`}
			jobID, err := d.AddJob(sleeperJob)
			Expect(err).To(BeNil())
			Expect(jobID).NotTo(Equal(""))

			err = d.Wait(jobID, time.Second*31, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			state, substate, err := d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(substate).To(Equal(""))
			// job was running through even there was exit code of 1
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))

			jobInfo, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())
			Expect(jobInfo.ID).To(Equal(jobID))
			Expect(jobInfo.ExitStatus).To(BeNumerically("==", 1))
			Expect(jobInfo.State.String()).To(Equal(drmaa2interface.Failed.String()))
		})

		It("should mark a job as failed when it can't be executed", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			sleeperJob.RemoteCommand = "/binary/notfound"
			jobID, err := d.AddJob(sleeperJob)
			Expect(err).To(BeNil())
			Expect(jobID).NotTo(Equal(""))

			err = d.Wait(jobID, time.Second*31, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			state, substate, err := d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(substate).To(Equal(""))
			Expect(state.String()).To(Equal(drmaa2interface.Failed.String()))

			jobInfo, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())
			Expect(jobInfo.ID).To(Equal(jobID))
			Expect(jobInfo.State.String()).To(Equal(drmaa2interface.Failed.String()))
			// jobInfo.SubState --
		})

	})

	Context("Job lifecyle management", func() {

		It("should wait until a job is running, suspend, resume and kill it", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			sleeperJob.Args = []string{"120"}

			jobID, err := d.AddJob(sleeperJob)
			Expect(err).To(BeNil())

			err = d.Wait(jobID, time.Second*31, drmaa2interface.Running)
			Expect(err).To(BeNil())

			err = d.JobControl(jobID, "suspend")
			Expect(err).To(BeNil())

			state, _, err := d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Suspended.String()))

			err = d.JobControl(jobID, "resume")
			Expect(err).To(BeNil())

			state, _, err = d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			err = d.JobControl(jobID, "terminate")
			Expect(err).To(BeNil())

			err = d.Wait(jobID, time.Second*40, drmaa2interface.Done, drmaa2interface.Failed, drmaa2interface.Undetermined)
			Expect(err).To(BeNil())

			state, _, err = d.JobState(jobID)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Failed.String()))
		})

	})

	Context("contact string", func() {

		It("should return the contact string of the drmaa connection", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()

			c, err := d.Contact()
			Expect(err).To(BeNil())
			Expect(c).NotTo(Equal(""))
		})

	})

	Measure("it should submit jobs in a short time", func(b Benchmarker) {
		<-time.Tick(time.Second * 5)

		d, err := NewDRMAATracker()
		Expect(err).To(BeNil())
		defer d.DestroySession()
		Expect(d).NotTo(BeNil())

		jobids := make([]string, 0, 16)
		submissiontime := b.Time("submissiontime", func() {
			jobid, _ := d.AddJob(sleeperJob)
			jobids = append(jobids, jobid)
		})

		Expect(submissiontime.Seconds()).To(BeNumerically("<", 0.050), "Submitting a job shouldn't take longer than 3 ms in avg.")

		// clean up
		for _, jobID := range jobids {
			d.JobControl(jobID, "terminate")
		}
		<-time.Tick(time.Second * 5)

	}, 20)

})
