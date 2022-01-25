package libdrmaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
)

var _ = Describe("Qhost", func() {

	Context("", func() {

		It("should return the machine list", func() {

			out := `HOSTNAME                ARCH         NCPU NSOC NCOR NTHR  LOAD  MEMTOT  MEMUSE  SWAPTO  SWAPUS
----------------------------------------------------------------------------------------------
global                  -               -    -    -    -     -       -       -       -       -
master                  lx-amd64        8    2    8    8  0.25    1.9G  334.9M 1024.0M  174.0M`

			hosts := ParseQhostForHostnames(out)
			Expect(len(hosts)).To(BeNumerically("==", 1))
			Expect(hosts[0].Name).To(Equal("master"))
			Expect(hosts[0].Architecture.String()).To(Equal(drmaa2interface.IA64.String()))
			Expect(hosts[0].OS.String()).To(Equal(drmaa2interface.Linux.String()))
			Expect(hosts[0].Sockets).To(BeNumerically("==", 2))
			Expect(hosts[0].CoresPerSocket).To(BeNumerically("==", 4))
			Expect(hosts[0].ThreadsPerCore).To(BeNumerically("==", 1))
			Expect(hosts[0].Load).To(BeNumerically("~", 0.24, 0.26))
		})

	})

})
