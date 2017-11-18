package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math"
)

var _ = Describe("Id", func() {

	Context("should return correct job ids", func() {

		It("should return 1 as first jobid", func() {
			SetJobID(0)
			jobid := GetNextJobID()
			Ω(jobid).To(Equal("1"))
			SetJobID(0)
		})

		It("should return 2 after 1", func() {
			SetJobID(0)
			jobid := GetNextJobID()
			Ω(jobid).To(Equal("1"))
			jobid = GetNextJobID()
			Ω(jobid).To(Equal("2"))
			SetJobID(0)
		})

		It("should do a rollover when max. number is reached", func() {
			SetJobID(math.MaxInt64)
			jobid := GetNextJobID()
			Ω(jobid).To(Equal("1"))
			SetJobID(0)
		})

	})

})
