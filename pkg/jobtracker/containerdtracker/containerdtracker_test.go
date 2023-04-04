package containerdtracker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
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
})
