package kubernetestracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("K8sState", func() {

	Context("Job state", func() {

		It("should return undetermined as state when job is not found", func() {
			cs, err := CreateClientSet()
			Ω(err).Should(BeNil())
			state := DRMAA2State(cs, "doesnotexist")
			Ω(state).Should(Equal(drmaa2interface.Undetermined))
		})

	})

})
