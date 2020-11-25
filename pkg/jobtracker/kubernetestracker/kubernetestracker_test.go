package kubernetestracker_test

import (
	"encoding/base64"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker"

	"os"
	"time"
)

var _ = Describe("KubernetesTracker", func() {

	Context("Basic interface test", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				//JobName:       "name1",
				RemoteCommand: "/bin/sh",
				JobCategory:   "busybox:latest",
				Args:          []string{"-c", "sleep 0"},
			}
			var err error
			kt, err = New("jobsession", nil)
			Ω(err).Should(BeNil())
		})

		WhenK8sIsAvailableIt("should be possible to AddJob()", func() {
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
		})

		WhenK8sIsAvailableIt("should be possible to DeleteJob()", func() {
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.DeleteJob(jobid)
			Ω(err).Should(BeNil())
		})

		WhenK8sIsAvailableIt("should be possible to AddArrayJob()", func() {
			jobid, err := kt.AddArrayJob(jt, 1, 2, 1, 0)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
		})

		WhenK8sIsAvailableIt("should be possible to ListJobs()", func() {
			jobids, err := kt.ListJobs()
			Ω(err).Should(BeNil())
			Ω(jobids).ShouldNot(BeNil())
		})

		WhenK8sIsAvailableIt("should be possible to ListArrayJobs()", func() {
			jobids, err := kt.ListArrayJobs("123")
			Ω(err).ShouldNot(BeNil())
			Ω(jobids).Should(BeNil())
		})

		WhenK8sIsAvailableIt("should be possible ListJobsCategories()", func() {
			cats, err := kt.ListJobCategories()
			Ω(err).Should(BeNil())
			Ω(cats).ShouldNot(BeNil())
			Ω(len(cats)).Should(BeNumerically("==", 0))
		})

	})

	Context("File staging", func() {

		It("should create map some content as a file inside the container", func() {
			jt := drmaa2interface.JobTemplate{
				JobCategory:   "busybox:latest",
				RemoteCommand: "cat",
				Args:          []string{"/my/file.txt", "/my/otherfile.txt"},
			}

			b64 := base64.StdEncoding.EncodeToString([]byte("content"))
			jt.StageInFiles = map[string]string{
				"configmap:" + b64: "/my/file.txt",
				"secret:" + b64:    "/my/otherfile.txt"}

			kt, err := New("jobsession", nil)
			Ω(err).Should(BeNil())

			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = kt.Wait(jobid, drmaa2interface.InfiniteTime,
				drmaa2interface.Done, drmaa2interface.Failed)
			Ω(err).Should(BeNil())
			state, _, err := kt.JobState(jobid)
			Ω(err).Should(BeNil())
			Ω(state.String()).Should(Equal(drmaa2interface.Done.String()))

			// delete configmap and secret
			err = kt.DeleteJob(jobid)
			Ω(err).Should(BeNil())
		})

	})

	Context("Unsupported interface functions", func() {
		var kt jobtracker.JobTracker

		BeforeEach(func() {
			var err error
			kt, err = New("", nil)
			Ω(err).Should(BeNil())
		})

		It("Unsupported ListJobCategories()", func() {
			_, err := kt.ListJobCategories()
			Ω(err).Should(BeNil())
		})

	})

	Context("JobSession related", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sh",
				JobCategory:   "busybox:latest",
			}
			var err error
			kt, err = New("jobsessionRelated", nil)
			Ω(err).Should(BeNil())
			// delete jobs from session if there are any remaining
			jobs, err := kt.ListJobs()
			Ω(err).Should(BeNil())
			for _, name := range jobs {
				Ω(kt.DeleteJob(name)).Should(BeNil())
			}
		})

		WhenK8sIsAvailableIt("ListJobs() should find the submitted jobs", func() {
			jt.Args = []string{"-c", "sleep 1"}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, err := kt.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 1))
			Ω(kt.DeleteJob(jobs[0])).Should(BeNil())
		})

	})

	Context("Basic Kubernetes Job Workflow", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				//JobName:       "workfloadtestjob",
				RemoteCommand: "/bin/sh",
				JobCategory:   "busybox:latest",
			}
			var err error
			kt, err = New("jobsession", nil)
			Ω(err).Should(BeNil())
		})

		WhenK8sIsAvailableIt("should be possible to track the states of a job life-cycle", func() {
			jt.Args = []string{"-c", "sleep 1"}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			Eventually(func() drmaa2interface.JobState {
				state, _, _ := kt.JobState(jobid)
				return state
			}, time.Second*30, time.Millisecond*50).Should(Equal(drmaa2interface.Running))

			Eventually(func() drmaa2interface.JobState {
				state, _, _ := kt.JobState(jobid)
				return state
			}, time.Second*30, time.Millisecond*250).Should(Equal(drmaa2interface.Done))
		})

		WhenK8sIsAvailableIt("should be possible to terminate a job", func() {
			jt.Args = []string{"-c", "sleep 10"}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			Eventually(func() drmaa2interface.JobState {
				state, _, _ := kt.JobState(jobid)
				return state
			}, time.Second*60, time.Millisecond*20).Should(Equal(drmaa2interface.Running))

			err = kt.JobControl(jobid, "terminate")
			Ω(err).Should(BeNil())

			Eventually(func() drmaa2interface.JobState {
				state, _, _ := kt.JobState(jobid)
				return state
			}, time.Second*60, time.Millisecond*10).Should(Equal(drmaa2interface.Undetermined))
		})

		WhenK8sIsAvailableIt("should be possible to wait for termination of a job", func() {
			jt.Args = []string{"-c", "sleep 10"}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			go func() {
				<-time.Tick(time.Millisecond * 333)
				kt.JobControl(jobid, "terminate")
			}()

			err = kt.Wait(jobid, time.Second*5, drmaa2interface.Failed, drmaa2interface.Undetermined)
			Ω(err).Should(BeNil())
			// TODO(DG) terminate should lead to failed state not undetermined
		})

		WhenK8sIsAvailableIt("should end in a failed state for a failing job", func() {
			jt.Args = []string{"-c", `exit 1`}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.Wait(jobid, time.Second*60, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Failed))
		})

		WhenK8sIsAvailableIt("should end in a done state for a successful job", func() {
			jt.Args = []string{"-c", `exit 0`}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.Wait(jobid, time.Second*60, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Done))
		})

		WhenK8sIsAvailableIt("should return JobInfo after the job is finished", func() {
			jt.Args = []string{"-c", `exit 0`}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.Wait(jobid, time.Second*60, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Done))
			ji, err := kt.JobInfo(jobid)
			Ω(err).Should(BeNil())
			Ω(ji.ID).Should(Equal(jobid))
			Ω(ji.State).Should(Equal(drmaa2interface.Done))
			Ω(ji.ExitStatus).Should(BeNumerically("==", 0))
		})

		WhenK8sIsAvailableIt("should return JobInfo after the job failed", func() {
			jt.Args = []string{"-c", `exit 1`}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.Wait(jobid, time.Second*60, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Failed))
			ji, err := kt.JobInfo(jobid)
			Ω(err).Should(BeNil())
			Ω(ji.ID).Should(Equal(jobid))
			Ω(ji.State).Should(Equal(drmaa2interface.Failed))
			Ω(ji.ExitStatus).Should(BeNumerically("==", 1))
		})

		WhenK8sIsAvailableIt("should finish the job when deadline is reached", func() {
			jt.Args = []string{"-c", "sleep 60"}
			jt.DeadlineTime = time.Now().Add(time.Second * 2)
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))
			err = kt.Wait(jobid, time.Second*30, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			Ω(kt.JobState(jobid)).Should(Equal(drmaa2interface.Failed))
		})

	})

	Context("Regression tests", func() {
		var kt jobtracker.JobTracker
		var jt drmaa2interface.JobTemplate

		BeforeEach(func() {
			jt = drmaa2interface.JobTemplate{
				//JobName:       "workfloadtestjob",
				RemoteCommand: "/bin/sh",
				JobCategory:   "busybox:latest",
			}
			var err error
			kt, err = New("jobsession", nil)
			Ω(err).Should(BeNil())
		})

		WhenK8sIsAvailableIt("should not crash when wait time is 0", func() {
			jt.Args = []string{"-c", `exit 0`}
			jobid, err := kt.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			kt.Wait(jobid, 0, drmaa2interface.Failed, drmaa2interface.Done, drmaa2interface.Undetermined)
		})

	})

	Context("Standard error cases", func() {
		WhenK8sIsAvailableIt("should fail to create a new tracker if k8s clientset can't be build", func() {
			home := os.Getenv("HOME")
			defer os.Setenv("HOME", home)
			os.Setenv("HOME", os.TempDir())
			track, err := New("", nil)
			Ω(err).ShouldNot(BeNil())
			Ω(track).Should(BeNil())
		})
	})

})
