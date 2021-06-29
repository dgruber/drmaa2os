package podmantracker_test

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
)

var _ = Describe("Podmantracker", func() {

	Context("e2e test", func() {

		var pt *PodmanTracker

		BeforeEach(func() {
			var err error

			if runtime.GOOS == "darwin" {
				// runs on macos connecting to a vagrant based vbox centos8 VM with
				// podman installed and remote socket activation with:
				// podman system service --time=0 unix:///tmp/podman.sock
				pt, err = New("testsession", PodmanTrackerParams{
					ConnectionURIOverride: "ssh://vagrant@localhost:2222/tmp/podman.sock?secure=False",
				})
				Expect(err).To(BeNil())
				Expect(pt).NotTo(BeNil())
			} else if runtime.GOOS == "linux" {
				_, err := exec.LookPath("podman")
				if err != nil {
					fmt.Printf("podman not in path")
					return
				}
				pt, err = New("testsession", PodmanTrackerParams{})
				Expect(err).To(BeNil())
				Expect(pt).NotTo(BeNil())
			}

		})

		It("should list images", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jc, err := pt.ListJobCategories()
			Expect(err).To(BeNil())
			Expect(jc).NotTo(BeNil())
			fmt.Printf("Container images: %v\n", jc)
		})

		It("should list containers", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jobs, err := pt.ListJobs()
			Expect(err).To(BeNil())
			Expect(jobs).NotTo(BeNil())
			fmt.Printf("Containers: %v\n", jobs)
		})

		It("should manage a complete container lifecylce", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jobid, err := pt.AddJob(drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "/bin/sleep",
				Args:          []string{"10"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			jobs, err := pt.ListJobs()
			Expect(err).To(BeNil())
			Expect(jobs).NotTo(BeNil())
			Expect(jobs).To(ContainElement(jobid))

			_, err = pt.JobInfo(jobid)
			Expect(err).To(BeNil())

			state, _, err := pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			err = pt.JobControl(jobid, "terminate")
			Expect(err).To(BeNil())

			err = pt.Wait(jobid, time.Second*120, drmaa2interface.Failed, drmaa2interface.Done)
			Expect(err).To(BeNil())

			state, _, err = pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Failed.String()))

			err = pt.DeleteJob(jobid)
			Expect(err).To(BeNil())
		})

		It("should wait until a successful job is finished", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jobid, err := pt.AddJob(drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			err = pt.Wait(jobid, time.Second*10, drmaa2interface.Done)
			Expect(err).To(BeNil())

			state, _, err := pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))
		})

		It("should wait until a failing job is finished", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jobid, err := pt.AddJob(drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "/bin/sh",
				Args:          []string{"-c", "exit 1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			err = pt.Wait(jobid, time.Second*10, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			state, _, err := pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Failed.String()))

			ji, err := pt.JobInfo(jobid)
			Expect(err).To(BeNil())
			Expect(ji.ExitStatus).To(BeNumerically("==", 1))
		})

		It("should wait timout while waiting", func() {
			if pt == nil {
				Skip("podman is not installed")
			}
			jobid, err := pt.AddJob(drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "/bin/sleep",
				Args:          []string{"5"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			err = pt.Wait(jobid, time.Second*1, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).NotTo(BeNil())

			state, _, err := pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))
		})

		// does not work with rootless containers and cgroups v1
		PIt("should suspend and resume the container", func() {
			if pt == nil {
				Skip("podman is not installed")
			}

			jobid, err := pt.AddJob(drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "/bin/sleep",
				Args:          []string{"5"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			state, _, err := pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			err = pt.JobControl(jobid, "suspend")
			Expect(err).To(BeNil())

			state, _, err = pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Suspended.String()))

			err = pt.JobControl(jobid, "resume")
			Expect(err).To(BeNil())

			state, _, err = pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			err = pt.Wait(jobid, time.Second*10, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).NotTo(BeNil())

			state, _, err = pt.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))

		})

	})

})
