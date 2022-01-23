package libdrmaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
)

var _ = Describe("Qhost", func() {

	Context("", func() {

		It("should return the machine list", func() {

			out := `HOSTNAME                ARCH         NCPU NSOC NCOR NTHR  LOAD  MEMTOT  MEMUSE  SWAPTO  SWAPUS
----------------------------------------------------------------------------------------------
global                  -               -    -    -    -     -       -       -       -       -
master                  lx-amd64        4    1    4    4  0.25    1.9G  334.9M 1024.0M  174.0M`

			hosts := ParseQhostForHostnames(out)
			Expect(len(hosts)).To(BeNumerically("==", 1))
			Expect(hosts[0]).To(Equal("master"))
		})

	})

})
