package drmaa2os_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dgruber/drmaa2os"
)

var _ = Describe("SessionmanagerHlp", func() {

	Context("reflect on ContactString in job tracker implementation specific params", func() {

		It("should set ContactString value", func() {
			type X struct {
				AnotherThing   int
				ContactString  string
				ContactString2 string
			}

			var test interface{} = X{
				AnotherThing:   11,
				ContactString:  "",
				ContactString2: "ö",
			}

			err := TryToSetContactString(&test, "ups")
			Expect(err).To(BeNil())
			Expect(test.(X).ContactString).To(Equal("ups"))
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
			Expect(test.(X).ContactString2).To(Equal("ö"))

			err = TryToSetContactString(&test, "fedex")
			Expect(err).To(BeNil())
			Expect(test.(X).ContactString).To(Equal("fedex"))
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
			Expect(test.(X).ContactString2).To(Equal("ö"))

			err = TryToSetContactString(&test, "cat")
			Expect(err).To(BeNil())
			Expect(test.(X).ContactString).To(Equal("cat"))
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
			Expect(test.(X).ContactString2).To(Equal("ö"))
		})

		It("should not crash if ContactString is not available", func() {
			type X struct {
				AnotherThing int
			}
			var test interface{} = X{
				AnotherThing: 11,
			}
			err := TryToSetContactString(&test, "∂")
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
			Expect(err).NotTo(BeNil())
			err = TryToSetContactString(nil, "k")
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
			Expect(err).NotTo(BeNil())
			err = TryToSetContactString(test, "k")
			Expect(err).NotTo(BeNil())
			Expect(test.(X).AnotherThing).To(BeNumerically("==", 11))
		})

		It("should set ContactString when struct is referenced as interface", func() {
			type X struct {
				ContactString string
			}
			test := X{
				ContactString: "",
			}
			err := TryToSetContactString(&test, "∂")
			Expect(err).NotTo(BeNil())
			Expect(test.ContactString).NotTo(Equal("∂"))

			f := func(asInterface interface{}) {
				err := TryToSetContactString(&asInterface, "†")
				Expect(err).To(BeNil())
				Expect(asInterface.(X).ContactString).To(Equal("†"))
			}
			f(test)
		})

	})

})
