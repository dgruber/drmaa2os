package drmaa2os_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dgruber/drmaa2os"
)

var _ = Describe("SessionmanagerHlp", func() {

	Context("reflect on ContactString in job tracker implementation specific params", func() {

		It("should set ContactString value", func() {
			test := struct {
				AnotherThing   int
				ContactString  string
				ContactString2 string
			}{
				AnotherThing:   11,
				ContactString:  "",
				ContactString2: "ö",
			}
			err := TryToSetContactString(&test, "ups")
			Expect(err).To(BeNil())
			Expect(test.ContactString).To(Equal("ups"))
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
			Expect(test.ContactString2).To(Equal("ö"))

			err = TryToSetContactString(&test, "fedex")
			Expect(err).To(BeNil())
			Expect(test.ContactString).To(Equal("fedex"))
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
			Expect(test.ContactString2).To(Equal("ö"))

			err = TryToSetContactString(&test, "cat")
			Expect(err).To(BeNil())
			Expect(test.ContactString).To(Equal("cat"))
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
			Expect(test.ContactString2).To(Equal("ö"))
		})

		It("should not crash if ContactString is not available", func() {
			test := struct {
				AnotherThing int
			}{
				AnotherThing: 11,
			}
			err := TryToSetContactString(&test, "∂")
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
			Expect(err).NotTo(BeNil())
			err = TryToSetContactString(nil, "k")
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
			Expect(err).NotTo(BeNil())
			err = TryToSetContactString(test, "k")
			Expect(err).NotTo(BeNil())
			Expect(test.AnotherThing).To(BeNumerically("==", 11))
		})

	})

})
