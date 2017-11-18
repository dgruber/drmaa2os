package helper_test

import (
	. "github.com/dgruber/drmaa2os/pkg/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {

	Context("Array Job ID convert functions", func() {

		It("should generate and resolve an array job ID into job IDs", func() {
			guids := []string{"1", "2", "3"}

			id := Guids2ArrayJobID(guids)
			guidsOut, err := ArrayJobID2GUIDs(id)

			Ω(err).Should(BeNil())
			Ω(guidsOut).Should(BeEquivalentTo(guids))
		})

	})

})
