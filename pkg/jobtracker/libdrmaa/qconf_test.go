package libdrmaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
)

var _ = Describe("Qconf", func() {

	Context("Basic tests", func() {

		It("should list all queues", func() {
			queues, err := QconfSQL()
			Expect(err).To(BeNil())
			Expect(len(queues)).To(BeNumerically(">=", 1))
			Expect(queues).To(ContainElement("all.q"))
		})

	})

})
