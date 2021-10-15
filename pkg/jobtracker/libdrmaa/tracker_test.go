package libdrmaa

import (
	"os"
	"time"

	"github.com/dgruber/drmaa2interface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tracker", func() {

	var sleeperJob drmaa2interface.JobTemplate

	getTempFile := func() string {
		file, _ := os.CreateTemp("", "drmaatracketest")
		name := file.Name()
		file.Close()
		return name
	}

	createTracker := func(standard bool) *DRMAATracker {
		if standard {
			standardTracker, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			return standardTracker
		}
		params := LibDRMAASessionParams{
			ContactString:           "",
			UsePersistentJobStorage: true,
			DBFilePath:              getTempFile(),
		}
		trackerWithParams, err := NewDRMAATrackerWithParams(params)
		Expect(err).To(BeNil())
		return trackerWithParams
	}

	BeforeEach(func() {
		sleeperJob = drmaa2interface.JobTemplate{
			RemoteCommand: "sleep",
			Args:          []string{"0"},
			JobName:       "mandatory",
		}
	})

	Context("Basic workflow", func() {

		It("should run a single job", func() {

			// run test for each tracker instantiation
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)
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

				d.DestroySession()
			}
		})

		It("should run a job array", func() {
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)
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

				d.DestroySession()
			}
		})

	})

	Context("Failing jobs", func() {

		It("should mark a job as done and show exit status when exit status != 0", func() {
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)
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

				d.DestroySession()
			}
		})

		It("should mark a job as failed when it can't be executed", func() {
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)
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

				d.DestroySession()
			}
		})

	})

	Context("Job lifecyle management", func() {

		It("should wait until a job is running, suspend, resume and kill it", func() {
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)
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

				d.DestroySession()
			}
		})

	})

	Context("contact string", func() {

		It("should return the contact string of the drmaa connection", func() {
			for _, standardTracker := range []bool{true, false} {
				d := createTracker(standardTracker)

				c, err := d.Contact()
				Expect(err).To(BeNil())
				Expect(c).NotTo(Equal(""))

				d.DestroySession()
			}
		})

	})

	Context("jobtracker with params", func() {

		It("should error when param has wrong type", func() {
			tracker, err := NewDRMAATrackerWithParams("string")
			Expect(err).NotTo(BeNil())
			Expect(tracker).To(BeNil())

			// this should fail as it must not be a reference
			tracker, err = NewDRMAATrackerWithParams(&LibDRMAASessionParams{})
			Expect(err).NotTo(BeNil())
			Expect(tracker).To(BeNil())
		})

		It("should error when param has wrong semantic", func() {
			tracker, err := NewDRMAATrackerWithParams(LibDRMAASessionParams{
				UsePersistentJobStorage: true,
				DBFilePath:              "",
			})
			Expect(err).NotTo(BeNil())
			Expect(tracker).To(BeNil())
		})

	})

	Measure("it should submit jobs in a short time", func(b Benchmarker) {
		<-time.Tick(time.Second * 5)

		for _, standardTracker := range []bool{true, false} {
			d := createTracker(standardTracker)
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

			d.DestroySession()
		}

	}, 20)

})
