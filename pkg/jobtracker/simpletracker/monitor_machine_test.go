package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MonitorHost", func() {

	It("should create a Machine struct of the local machine", func() {
		machine, err := GetLocalMachineInfo()
		Ω(err).Should(BeNil())
		Ω(machine.Name).ShouldNot(Equal(""))
	})

})
