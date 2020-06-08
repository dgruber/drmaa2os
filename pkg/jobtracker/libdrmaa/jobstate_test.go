package libdrmaa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Jobstate", func() {

	Context("basic tests", func() {

		It("should convert the state", func() {
			state := ConvertDRMAAStateToDRMAA2State(drmaa.PsUndetermined)
			Expect(state).To(Equal(drmaa2interface.Undetermined))

			Expect(ConvertDRMAAStateToDRMAA2State(drmaa.PsQueuedActive)).To(Equal(drmaa2interface.Queued))
			Expect(ConvertDRMAAStateToDRMAA2State(drmaa.PsRunning)).To(Equal(drmaa2interface.Running))
			Expect(ConvertDRMAAStateToDRMAA2State(drmaa.PsDone)).To(Equal(drmaa2interface.Done))
			Expect(ConvertDRMAAStateToDRMAA2State(drmaa.PsFailed)).To(Equal(drmaa2interface.Failed))
		})
	})

})
