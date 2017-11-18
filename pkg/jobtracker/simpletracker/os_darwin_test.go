package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os/exec"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("OsDarwin", func() {

	Context("basic tests", func() {
		It("should return the job state for a sleeper", func() {

			cmd := exec.Command("/bin/sleep", "1")
			cmd.Start()
			pid := cmd.Process.Pid

			state, err := OSStateStringforPID(fmt.Sprintf("%d", pid))
			Ω(err).Should(BeNil())
			Ω(state).Should(ContainSubstring("S"))
		})

		It("should indentify a suspended process", func() {
			Ω(OSStateToDRMAA2State("T")).Should(Equal(drmaa2interface.Suspended))
		})

		It("should indentify a running process", func() {
			Ω(OSStateToDRMAA2State("S")).Should(Equal(drmaa2interface.Running))
		})
	})

})
