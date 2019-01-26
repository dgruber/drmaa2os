package slurmcli_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/slurmcli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Tracker", func() {

	var s *Slurm

	BeforeEach(func() {
		s = NewSlurm("./fakes/sbatch.sh",
			"./fakes/squeue.sh",
			"./fakes/scontrol.sh",
			"./fakes/scancel.sh",
			"./fakes/sacct.sh",
			false)
	})

	Context("Basic operations", func() {

		It("should be possible to create a job tracker", func() {
			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
		})

		It("should be possible to submit a job", func() {
			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "mycommand",
			})
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
		})

		It("should be possible to list jobs", func() {
			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
			jobs, err := tracker.ListJobs()
			Ω(err).Should(BeNil())
			Ω(jobs).ShouldNot(BeNil())
		})

	})

	Context("Job state related", func() {

		It("should detect all states", func() {
			s = NewSlurm("./fakes/sbatch.sh",
				"./fakes/squeue.sh",
				"./fakes/scontrol.sh",
				"./fakes/scancel.sh",
				"./fakes/state.sh",
				false)

			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			// fakes/state.sh file contains the mapping
			Ω(tracker.JobState("1")).Should(Equal(drmaa2interface.Running))
			Ω(tracker.JobState("2")).Should(Equal(drmaa2interface.Running))
			Ω(tracker.JobState("3")).Should(Equal(drmaa2interface.Done))
			Ω(tracker.JobState("4")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("5")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("6")).Should(Equal(drmaa2interface.Running))
			Ω(tracker.JobState("7")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("8")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("9")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("10")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("11")).Should(Equal(drmaa2interface.Queued))
			Ω(tracker.JobState("12")).Should(Equal(drmaa2interface.Suspended))
			Ω(tracker.JobState("13")).Should(Equal(drmaa2interface.QueuedHeld))
			Ω(tracker.JobState("14")).Should(Equal(drmaa2interface.Queued))
			Ω(tracker.JobState("15")).Should(Equal(drmaa2interface.RequeuedHeld))
			Ω(tracker.JobState("16")).Should(Equal(drmaa2interface.Requeued))
			Ω(tracker.JobState("17")).Should(Equal(drmaa2interface.Running))
			Ω(tracker.JobState("18")).Should(Equal(drmaa2interface.Undetermined))
			Ω(tracker.JobState("19")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("20")).Should(Equal(drmaa2interface.Requeued))
			Ω(tracker.JobState("21")).Should(Equal(drmaa2interface.Running))
			Ω(tracker.JobState("22")).Should(Equal(drmaa2interface.Suspended))
			Ω(tracker.JobState("23")).Should(Equal(drmaa2interface.Suspended))
			Ω(tracker.JobState("24")).Should(Equal(drmaa2interface.Failed))
			Ω(tracker.JobState("99")).Should(Equal(drmaa2interface.Undetermined))
		})

	})

	Context("Error cases", func() {

		It("should fail to create a job tracker when commands are not available", func() {
			commands := []string{
				"./fakes/sbatch.sh",
				"./fakes/squeue.sh",
				"./fakes/scontrol.sh",
				"./fakes/scancel.sh",
				"./fakes/sacct.sh"}

			for i := 0; i < 5; i++ {
				orig := commands[i]
				commands[i] = "notfound.sh"
				s = NewSlurm(commands[0], commands[1], commands[2],
					commands[3], commands[4], true)
				tracker, err := New("test", s)
				Ω(err).ShouldNot(BeNil())
				Ω(tracker).Should(BeNil())
				commands[i] = orig
			}
		})

		It("should fail to run a job when RemoteCommand is not set", func() {
			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{})
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

		It("should fail to run a job when job submission fails is not set", func() {
			tracker, err := New("test", s)
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
			jobid, err := tracker.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "fail",
			})
			Ω(err).ShouldNot(BeNil())
			Ω(jobid).Should(Equal(""))
		})

	})

})
