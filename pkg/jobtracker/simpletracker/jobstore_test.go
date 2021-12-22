package simpletracker_test

import (
	"os"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Jobstore", func() {

	Context("Basic JobStore operations", func() {

		var inmemory *JobStore
		var persistent *PersistentJobStorage

		BeforeEach(func() {
			inmemory = NewJobStore()
			Ω(inmemory).ShouldNot(BeNil())

			file, err := os.CreateTemp("", "jobstoretest")
			Expect(err).To(BeNil())
			name := file.Name()
			file.Close()
			persistent, err = NewPersistentJobStore(name)
			Expect(err).To(BeNil())
		})

		It("should be possible to create a JobStore, save a job, and get the PID and jobTemplate", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				Ω(store).ShouldNot(BeNil())
				store.SaveJob("13", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 77)
				store.SaveJob("1", drmaa2interface.JobTemplate{RemoteCommand: "rc2"}, 13)
				store.SaveJob("12", drmaa2interface.JobTemplate{RemoteCommand: "rc3"}, 10)
				pid, err := store.GetPID("12")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 10))
				pid, err = store.GetPID("1")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 13))
				pid, err = store.GetPID("13")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 77))
				jt, err := store.GetJobTemplate("12")
				Ω(err).Should(BeNil())
				Ω(jt.RemoteCommand).Should(Equal("rc3"))
				jt, err = store.GetJobTemplate("1")
				Ω(err).Should(BeNil())
				Ω(jt.RemoteCommand).Should(Equal("rc2"))
				jt, err = store.GetJobTemplate("13")
				Ω(err).Should(BeNil())
				Ω(jt.RemoteCommand).Should(Equal("rc"))
			}
		})

		It("should find PID of array job task", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveArrayJob("13",
					[]int{77, 78, 79},
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				store.SaveJob("13.1", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 77)
				store.SaveJob("13.2", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 78)
				store.SaveJob("13.3", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 79)
				pid, err := store.GetPID("13.1")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 77))
				pid, err = store.GetPID("13.2")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 78))
				pid, err = store.GetPID("13.3")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 79))
			}
		})

		It("should error when job is not found", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				pid, err := store.GetPID("12")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
			}
		})

		It("should error when job id is wrong", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				pid, err := store.GetPID("12.asdf")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
				pid, err = store.GetPID("..")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
				store.SaveJob("13.2", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 77)
				pid, err = store.GetPID("13.asdf")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
			}
		})

		It("should error when task is not found", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveJob("13.1", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 77)
				pid, err := store.GetPID("13.77")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
				pid, err = store.GetPID("13.abc")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
			}
		})

		It("should error when task is not found", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveArrayJob("13",
					[]int{77, 78, 79},
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				pid, err := store.GetPID("13.10")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
				pid, err = store.GetPID("13.abc")
				Ω(err).ShouldNot(BeNil())
				Ω(pid).Should(BeNumerically("==", -1))
			}
		})

		It("should save and delete a job array", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveJob("77.2", drmaa2interface.JobTemplate{RemoteCommand: "rc"}, 77)
				store.SaveArrayJob("13",
					[]int{77, 78, 79},
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				Ω(store.HasJob("13.1")).Should(BeTrue())
				Ω(store.HasJob("13.2")).Should(BeTrue())
				Ω(store.HasJob("13.3")).Should(BeTrue())
				store.RemoveJob("13.2")
				Ω(store.HasJob("13.2")).Should(BeFalse())
				store.RemoveJob("13")
				Ω(store.HasJob("13.1")).Should(BeFalse())
				Ω(store.HasJob("13.3")).Should(BeFalse())
			}
		})

		It("should save and a job array and add the PID of a task afterwards", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveArrayJob("13",
					[]int{0, 0, 0},
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				pid, err := store.GetPID("13.2")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 0))

				err = store.SaveArrayJobPID("13", 2, 77)
				Ω(err).Should(BeNil())
				pid, err = store.GetPID("13.2")
				Ω(err).Should(BeNil())
				Ω(pid).Should(BeNumerically("==", 77))

				err = store.SaveArrayJobPID("13", 50, 77)
				Ω(err).ShouldNot(BeNil())

				err = store.SaveArrayJobPID("50", 50, 77)
				Ω(err).ShouldNot(BeNil())
			}
		})

		It("should return the task IDs of an array job", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveArrayJob("112",
					[]int{0, 0, 0}, // no pid
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				tasks := store.GetArrayJobTaskIDs("112")
				Ω(len(tasks)).To(BeNumerically("==", 3))
				Ω(tasks[0]).To(Equal("112.1"))
				Ω(tasks[1]).To(Equal("112.2"))
				Ω(tasks[2]).To(Equal("112.3"))
			}
		})

		It("should return the task IDs of an array job as job IDs", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				store.SaveArrayJob("112",
					[]int{0, 0, 0}, // no pid
					drmaa2interface.JobTemplate{RemoteCommand: "rc"},
					1, 3, 1)
				tasks := store.GetJobIDs()
				Ω(len(tasks)).To(BeNumerically("==", 3))
				Ω(tasks).To(ContainElement("112.1"))
				Ω(tasks).To(ContainElement("112.2"))
				Ω(tasks).To(ContainElement("112.3"))
			}
		})

		It("should store the job info and return it", func() {
			for _, store := range []JobStorer{persistent, inmemory} {
				ji := drmaa2interface.JobInfo{
					ID:         "id",
					ExitStatus: 1,
					JobOwner:   "owner",
				}
				ji.ExtensionList = map[string]string{
					"extension": "value",
				}
				err := store.SaveJobInfo("id", ji)
				Expect(err).Should(BeNil())
				jiBack, err := store.GetJobInfo("id")
				Expect(err).Should(BeNil())
				Expect(jiBack.ID).Should(Equal("id"))
				Expect(jiBack.JobOwner).Should(Equal("owner"))
				Expect(jiBack.ExitStatus).Should(BeNumerically("==", 1))
				Expect(jiBack.ExtensionList).ShouldNot(BeNil())
				Expect(jiBack.ExtensionList["extension"]).Should(Equal("value"))
			}
		})

	})

	Context("Persistent JobStore operations", func() {

		It("should fail to create a persistent job storage when DB file is not set", func() {
			var err error
			persistent, err := NewPersistentJobStore("")
			Expect(err).NotTo(BeNil())
			Expect(persistent).To(BeNil())
		})

	})

})
