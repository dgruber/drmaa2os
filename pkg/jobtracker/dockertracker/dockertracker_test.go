package dockertracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgruber/drmaa2interface"
)

var _ = BeforeSuite(func() {
	// pull required images

	// pull alpine
	st := simpletracker.New("")
	jobID, err := st.AddJob(drmaa2interface.JobTemplate{
		RemoteCommand: "docker",
		Args:          []string{"pull", "alpine"},
	})
	Ω(err).Should(BeNil())
	err = st.Wait(jobID, time.Second*120, drmaa2interface.Done)
	Ω(err).Should(BeNil())
})

var _ = Describe("Dockertracker", func() {

	Context("Creation and destruction", func() {

		It("should be possible to create a tracker when docker is available", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())
		})

	})

	Context("List container images as Job Classes", func() {

		It("should list without errors", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			images, err := tracker.ListJobCategories()
			Ω(err).Should(BeNil())
			Ω(len(images)).Should(BeNumerically(">=", 0))
		})

		It("should throw an error when tracker was not initialized", func() {
			var tracker DockerTracker
			images, err := tracker.ListJobCategories()

			Ω(err).ShouldNot(BeNil())
			Ω(images).Should(BeNil())
		})

	})

	Context("Add job", func() {

		var jt drmaa2interface.JobTemplate
		jt.ExtensionList = map[string]string{"exposedPorts": "8080/tcp"}

		var tracker *DockerTracker

		BeforeEach(func() {
			tracker, _ = New("")
			jt = drmaa2interface.JobTemplate{
				RemoteCommand:  "/bin/sleep",
				Args:           []string{"1"},
				JobCategory:    "alpine",
				StageInFiles:   map[string]string{"README.md": "/README.md"},
				JobEnvironment: map[string]string{"test": "value"},
			}
		})

		It("should add the job without error", func() {
			id, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(id).ShouldNot(Equal(""))

			state, _, _ := tracker.JobState(id)
			Ω(err).Should(BeNil())
			Ω(state).Should(Equal(drmaa2interface.Running))
		})

		It("should fail adding the job when JobCategory in job template is missing", func() {
			jt.JobCategory = ""
			id, err := tracker.AddJob(jt)
			Ω(err).ShouldNot(BeNil())
			Ω(id).Should(Equal(""))
		})

		It("should print output to file", func() {
			os.Remove("./testfile")

			jt.RemoteCommand = "/bin/sh"
			jt.Args = []string{"-c", `echo prost`}
			jt.OutputPath = "./testfile"

			id, err := tracker.AddJob(jt)

			Ω(err).Should(BeNil())
			Ω(id).ShouldNot(Equal(""))
			err = tracker.Wait(id, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())
			content, err := ioutil.ReadFile("./testfile")
			Ω(err).Should(BeNil())
			Ω(string(content)).Should(ContainSubstring("prost"))
			os.Remove("./testfile")
		})

		It("should print stderr to file", func() {
			jt.RemoteCommand = "/bin/sh"
			jt.Args = []string{"-c", `date illegal`}
			jt.ErrorPath = "./errtestfile"

			id, err := tracker.AddJob(jt)

			Ω(err).Should(BeNil())
			Ω(id).ShouldNot(Equal(""))
			err = tracker.Wait(id, 5*time.Second, drmaa2interface.Done, drmaa2interface.Failed)
			Ω(err).Should(BeNil())
			content, err := ioutil.ReadFile("./errtestfile")
			Ω(err).Should(BeNil())
			Ω(string(content)).Should(ContainSubstring("date: invalid date"))
			os.Remove("./errtestfile")
		})
	})

	Context("Array job", func() {

		jt := drmaa2interface.JobTemplate{
			RemoteCommand: "/bin/sleep",
			Args:          []string{"1"},
			JobCategory:   "alpine",
		}

		var tracker *DockerTracker

		BeforeEach(func() { tracker, _ = New("") })

		It("should add the job without error", func() {
			ids, err := tracker.AddArrayJob(jt, 1, 10, 1, 0)
			Ω(err).Should(BeNil())
			Ω(ids).ShouldNot(Equal(""))

			jobids, err := tracker.ListArrayJobs(ids)
			Ω(err).Should(BeNil())
			Ω(jobids).ShouldNot(BeNil())
			Ω(len(jobids)).Should(BeNumerically("==", 10))
		})

	})

	Context("Job life cycle", func() {

		jt := drmaa2interface.JobTemplate{
			RemoteCommand: "/bin/sleep",
			Args:          []string{"1"},
			JobCategory:   "alpine",
		}

		var tracker *DockerTracker

		BeforeEach(func() { tracker, _ = New("") })

		It("add job and wait until finished", func() {
			id, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(id).ShouldNot(Equal(""))

			state, _, _ := tracker.JobState(id)
			Ω(err).Should(BeNil())
			Ω(state.String()).Should(Equal(drmaa2interface.Running.String()))

			err = tracker.Wait(id, drmaa2interface.InfiniteTime, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			state, _, _ = tracker.JobState(id)
			Ω(err).Should(BeNil())
			Ω(state.String()).Should(Equal(drmaa2interface.Done.String()))
		})

		It("add job and terminate", func() {
			id, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(id).ShouldNot(Equal(""))

			state, _, _ := tracker.JobState(id)
			Ω(err).Should(BeNil())
			Ω(state.String()).Should(Equal(drmaa2interface.Running.String()))

			err = tracker.JobControl(id, "terminate")
			Ω(err).Should(BeNil())

			state, _, _ = tracker.JobState(id)
			Ω(state.String()).Should(Equal(drmaa2interface.Failed.String()))

			fmt.Println(id)
			err = tracker.DeleteJob(id)
			Ω(err).Should(BeNil())

			state, _, _ = tracker.JobState(id)
			Ω(state.String()).Should(Equal(drmaa2interface.Undetermined.String()))
		})

	})

	Context("List containers as Jobs", func() {

		It("should list without errors", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			jobs, err := tracker.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically(">=", 0))
		})

		It("should throw an error when tracker was not initialized", func() {
			var tracker DockerTracker
			jobs, err := tracker.ListJobs()

			Ω(err).ShouldNot(BeNil())
			Ω(jobs).Should(BeNil())
		})

		It("should list jobs from the job session", func() {
			tracker, err := New("testsessionXY")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			jobs, err := tracker.ListJobs()
			Ω(err).Should(BeNil())
			// delete all remaining jobs
			for _, id := range jobs {
				tracker.DeleteJob(id)
			}

			// alpine image needs to be pulled before!
			jt := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
				JobCategory:   "alpine",
			}

			jobid, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			jobs, err = tracker.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 1))

			err = tracker.Wait(jobid, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			jobs, err = tracker.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 1))

			err = tracker.DeleteJob(jobid)
			Ω(err).Should(BeNil())

			jobs, err = tracker.ListJobs()
			Ω(err).Should(BeNil())
			Ω(len(jobs)).Should(BeNumerically("==", 0))
		})

	})

	Context("Job template", func() {

		It("should return the job template of a job", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			jt := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"0"},
				JobCategory:   "alpine",
				OutputPath:    "/dev/stdout",
				ErrorPath:     "/dev/stderr",
				Extension: drmaa2interface.Extension{
					ExtensionList: map[string]string{
						"some": "extension",
					},
				},
			}

			jobid, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())

			jt2, err := tracker.JobTemplate(jobid)
			Ω(err).Should(BeNil())

			Ω(jt2.RemoteCommand).Should(Equal(jt.RemoteCommand))
			Ω(jt2.Args).Should(Equal(jt.Args))
			Ω(jt2.JobCategory).Should(Equal(jt.JobCategory))
			Ω(jt2.OutputPath).Should(Equal(jt.OutputPath))
			Ω(jt2.ErrorPath).Should(Equal(jt.ErrorPath))
			Ω(jt2.Extension.ExtensionList["some"]).Should(Equal("extension"))

			tracker.Wait(jobid, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			err = tracker.DeleteJob(jobid)
			Ω(err).Should(BeNil())
		})

	})

	Context("Job output", func() {

		It("should return the standard output of a job", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			// create temporary file
			tmpFile, err := ioutil.TempFile("", "drmaa2os")
			Ω(err).Should(BeNil())
			Ω(tmpFile).ShouldNot(BeNil())
			tmpFile.Close()

			defer os.Remove(tmpFile.Name())

			jt := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/echo",
				Args:          []string{"test"},
				JobCategory:   "alpine",
				OutputPath:    tmpFile.Name(),
			}

			jobid, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			output, err := ioutil.ReadFile(tmpFile.Name())
			Ω(err).Should(BeNil())
			Ω(string(output)).Should(Equal("test\n"))
		})

		It("should return the error output of a job", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			// create temporary file
			tmpFile, err := ioutil.TempFile("", "drmaa2os")
			Ω(err).Should(BeNil())
			Ω(tmpFile).ShouldNot(BeNil())
			tmpFile.Close()

			defer os.Remove(tmpFile.Name())

			jt := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sh",
				Args:          []string{"-c", "echo test 1>&2"},
				JobCategory:   "alpine",
				ErrorPath:     tmpFile.Name(),
			}

			jobid, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			output, err := ioutil.ReadFile(tmpFile.Name())
			Ω(err).Should(BeNil())
			Ω(string(output)).Should(Equal("test\n"))
		})

		It("should return the standard and error output of a job", func() {
			tracker, err := New("")
			Ω(err).Should(BeNil())
			Ω(tracker).ShouldNot(BeNil())

			// create temporary file 1
			tmpFile, err := ioutil.TempFile("", "drmaa2os")
			Ω(err).Should(BeNil())
			Ω(tmpFile).ShouldNot(BeNil())
			tmpFile.Close()

			defer os.Remove(tmpFile.Name())

			// create temporary file 2
			tmpFile2, err := ioutil.TempFile("", "drmaa2os")
			Ω(err).Should(BeNil())
			Ω(tmpFile2).ShouldNot(BeNil())
			tmpFile2.Close()

			defer os.Remove(tmpFile2.Name())

			jt := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sh",
				Args:          []string{"-c", "echo testtest 1>&2 && echo test"},
				JobCategory:   "alpine",
				OutputPath:    tmpFile.Name(),
				ErrorPath:     tmpFile2.Name(),
			}

			jobid, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())
			Ω(jobid).ShouldNot(Equal(""))

			err = tracker.Wait(jobid, drmaa2interface.InfiniteTime, drmaa2interface.Done)
			Ω(err).Should(BeNil())

			output, err := ioutil.ReadFile(tmpFile.Name())
			Ω(err).Should(BeNil())
			Ω(string(output)).Should(Equal("test\n"))

			output, err = ioutil.ReadFile(tmpFile2.Name())
			Ω(err).Should(BeNil())
			Ω(string(output)).Should(Equal("testtest\n"))
		})

	})

})
