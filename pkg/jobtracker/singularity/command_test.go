package singularity

import (
	"github.com/dgruber/drmaa2interface"
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = g.Describe("Command", func() {

	var jt = drmaa2interface.JobTemplate{}

	g.BeforeEach(func() {
		jt = drmaa2interface.JobTemplate{
			RemoteCommand: "run",
			Args:          []string{"arg1", "arg2"},
			JobCategory:   "image",
		}
	})

	g.Context("Extensions", func() {
		extension := func(name, value string) (string, []string) {
			jt.ExtensionList = map[string]string{
				name: value,
			}
			return createCommandAndArgs(jt)
		}
		g.It("should insert boolean extensions", func() {
			cmd, args := extension("pid", "")
			Ω(cmd).Should(Equal("singularity"))
			Ω(args).Should(ContainElement("--pid"))
			Ω(args).ShouldNot(ContainElement(""))

			cmd, args = extension("hostname", "alleswurst")
			Ω(cmd).Should(Equal("singularity"))
			Ω(args).Should(ContainElement("--hostname"))
			Ω(args).Should(ContainElement("alleswurst"))
		})
	})
})
