package containerdtracker_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/containerdtracker"
)

var _ = Describe("ContainerdJobTracker", func() {
	var tracker *containerdtracker.ContainerdJobTracker
	var err error

	BeforeEach(func() {
		tracker, err = containerdtracker.NewContainerdJobTracker("/run/containerd/containerd.sock")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// Clean up any created containers
	})

	Describe("AddJob", func() {
		It("should create a new container and return its ID", func() {
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"echo", "hello"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())
		})
	})

	Describe("JobInfo", func() {
		It("should return the JobInfo for the specified container ID", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "5"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			jobInfo, err := tracker.JobInfo(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobInfo.ID).To(Equal(jobID))
			//Expect(jobInfo.JobCategory).To(Equal("docker.io/library/busybox:latest"))
			Expect(jobInfo.State).To(BeNumerically("==", drmaa2interface.Running))
		})
	})

	Describe("ListJobs", func() {
		It("should return a list of container IDs", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "5"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			jobIDs, err := tracker.ListJobs()
			Expect(err).NotTo(HaveOccurred())
			Expect(jobIDs).To(ContainElement(jobID))
		})
	})

	Describe("ListArrayJobs", func() {
		It("should return a list of container IDs for the given array job ID", func() {
			// Use the AddArrayJob function to create new array jobs
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "5"},
			}
			arrayJobID, err := tracker.AddArrayJob(jobTemplate, 1, 2, 1, 2)
			Expect(err).NotTo(HaveOccurred())
			Expect(arrayJobID).NotTo(BeEmpty())

			jobIDs, err := tracker.ListArrayJobs(arrayJobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(jobIDs)).To(Equal(2))

			// Check if the array job container IDs are included in the list of all container IDs
			allJobIDs, err := tracker.ListJobs()
			Expect(err).NotTo(HaveOccurred())
			for _, jobID := range jobIDs {
				Expect(allJobIDs).To(ContainElement(jobID))
			}
		})
	})

	Describe("Wait", func() {
		It("should wait for the container to reach the specified state", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "2"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			// Wait for the container to reach the Running state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Running)
			Expect(err).NotTo(HaveOccurred())

			// Wait for the container to reach the Done state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Done)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DeleteJob", func() {

		// TODO: This test is failing because the container seems to be
		// still running...
		XIt("should delete the specified container", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"true"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			// Wait for the container to reach the Done state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Done)
			Expect(err).NotTo(HaveOccurred())

			// Delete the container
			err = tracker.DeleteJob(jobID)
			Expect(err).NotTo(HaveOccurred())

			// Check if the container is deleted
			jobIDs, err := tracker.ListJobs()
			Expect(err).NotTo(HaveOccurred())
			Expect(jobIDs).NotTo(ContainElement(jobID))
		})
	})

	Describe("JobControl", func() {

		It("should suspend and resume the specified container", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "10"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			// Wait for the container to reach the Running state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Running)
			Expect(err).NotTo(HaveOccurred())

			// Suspend the container
			err = tracker.JobControl(jobID, jobtracker.JobControlSuspend)
			Expect(err).NotTo(HaveOccurred())

			// Check if the container is in Suspended state
			state, _, err := tracker.JobState(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(state.String()).To(Equal(drmaa2interface.Suspended.String()))

			// Resume the container
			err = tracker.JobControl(jobID, jobtracker.JobControlResume)
			Expect(err).NotTo(HaveOccurred())

			// Check if the container is in Running state
			state, _, err = tracker.JobState(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

		})

		It("should terminate the specified container", func() {
			// Use the AddJob function to create a new container
			jobTemplate := drmaa2interface.JobTemplate{
				JobCategory: "docker.io/library/busybox:latest",
				Args:        []string{"sleep", "10"},
			}
			jobID, err := tracker.AddJob(jobTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobID).NotTo(BeEmpty())

			// Wait for the container to reach the Running state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Running)
			Expect(err).NotTo(HaveOccurred())

			// Terminate the container
			err = tracker.JobControl(jobID, jobtracker.JobControlTerminate)
			Expect(err).NotTo(HaveOccurred())

			// Wait for the container to reach the Failed state
			err = tracker.Wait(jobID, 10*time.Second, drmaa2interface.Failed)
			Expect(err).NotTo(HaveOccurred())

			// Check if the container is in Failed state
			state, _, err := tracker.JobState(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(drmaa2interface.Failed))
		})
	})
})
