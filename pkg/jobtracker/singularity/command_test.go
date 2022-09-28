package singularity

import (
	"github.com/dgruber/drmaa2interface"
	g "github.com/onsi/ginkgo/v2"
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
			Ω(args).Should(BeEquivalentTo([]string{"exec", "--pid", "image", "run", "arg1", "arg2"}))

			cmd, args = extension("hostname", "alleswurst")
			Ω(cmd).Should(Equal("singularity"))
			Ω(args).Should(ContainElement("--hostname"))
			Ω(args).Should(ContainElement("alleswurst"))
		})
		g.It("should insert global arguments from the extensions", func() {
			jt.ExtensionList = map[string]string{
				"pid":   "",
				"debug": "true",
			}
			cmd, args := createCommandAndArgs(jt)
			Ω(cmd).Should(Equal("singularity"))
			Ω(args).Should(BeEquivalentTo([]string{"--debug", "exec", "--pid", "image", "run", "arg1", "arg2"}))
		})

	})
})
