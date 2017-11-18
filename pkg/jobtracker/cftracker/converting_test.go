package cftracker

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/cftracker/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Converting", func() {

	Context("convert tasks in names", func() {
		It("should never return nil", func() {
			names := convertTasksInNames(nil)
			Ω(names).ShouldNot(BeNil())
		})
	})

	Context("convert task in JobInfo", func() {
		It("should convert a failed task", func() {
			jobInfo := convertTaskInJobinfo(*fake.FailedTaskFake())
			Ω(jobInfo.State).Should(Equal(drmaa2interface.Failed))
		})
		It("should convert a succeeded task", func() {
			jobInfo := convertTaskInJobinfo(*fake.SucceededTaskFake())
			Ω(jobInfo.State).Should(Equal(drmaa2interface.Done))
		})
	})

	Context("convert JobTemplate in task requests", func() {

		jt := drmaa2interface.JobTemplate{
			RemoteCommand:     "/bin/sleep",
			Args:              []string{"123"},
			JobCategory:       "123-123-123-123",
			WorkingDirectory:  "/working/dir",
			CandidateMachines: []string{"hostname"},
			MinPhysMemory:     512, // bytes
		}

		It("should convert a complete JobTemplate correctly", func() {
			tr, err := convertJobTemplateInTaskRequest(jt)
			Ω(err).Should(BeNil())
			Ω(tr.Command).Should(Equal("/bin/sleep 123"))
			Ω(tr.DropletGUID).Should(Equal("123-123-123-123"))
			Ω(tr.MemoryInMegabyte).Should(BeNumerically("==", 1))
		})
	})

})
