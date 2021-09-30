package storage_test

import (
	. "github.com/dgruber/drmaa2os/pkg/storage"
	. "github.com/dgruber/drmaa2os/pkg/storage/boltstore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

var _ = Describe("Store", func() {
	var store Storer

	BeforeEach(func() {
		os.Remove("test.db")
		store = NewBoltStore("test.db")
	})

	Describe("Calling init function", func() {
		It("should have no error", func() {
			err := store.Init()
			Expect(err).To(BeNil())
			store.Exit()
		})
	})

	Describe("Put and get values", func() {
		Context("and working with just one value", func() {
			It("should get the same value back as we put in", func() {
				errInit := store.Init()
				Expect(errInit).To(BeNil())
				defer store.Exit()
				errPut := store.Put(JobSessionType, "key", "value")
				Expect(errPut).To(BeNil())
				value, errGet := store.Get(JobSessionType, "key")
				Expect(errGet).To(BeNil())
				Expect(value).To(BeIdenticalTo("value"))
			})
		})
		Context("and working with multiple values in multiple buckets", func() {
			It("should get the same values back for each bucket", func() {
				errInit := store.Init()
				Expect(errInit).To(BeNil())
				defer store.Exit()
				errPut := store.Put(JobSessionType, "key1", "value1")
				Expect(errPut).To(BeNil())
				errPut2 := store.Put(ReservationSessionType, "key2", "value2")
				Expect(errPut2).To(BeNil())
				value, errGet := store.Get(JobSessionType, "key1")
				Expect(errGet).To(BeNil())
				Expect(value).To(BeIdenticalTo("value1"))
				value2, errGet2 := store.Get(ReservationSessionType, "key2")
				Expect(errGet2).To(BeNil())
				Expect(value2).To(BeIdenticalTo("value2"))
			})
		})
	})

	Describe("Testing List() functionality", func() {
		It("should list no keys after creation", func() {
			errInit := store.Init()
			Expect(errInit).To(BeNil())
			defer store.Exit()
			list, errList := store.List(JobSessionType)
			Expect(errList).To(BeNil())
			Ω(list).ShouldNot(BeNil())
			Ω(len(list)).Should(BeNumerically("==", 0))
		})
		It("should list all keys", func() {
			errInit := store.Init()
			Expect(errInit).To(BeNil())
			defer store.Exit()
			errPut1 := store.Put(JobSessionType, "key1", "value1")
			Expect(errPut1).To(BeNil())
			errPut2 := store.Put(JobSessionType, "key2", "value2")
			Expect(errPut2).To(BeNil())
			errPut3 := store.Put(JobSessionType, "key3", "value3")
			Expect(errPut3).To(BeNil())
			list, errList := store.List(JobSessionType)
			Expect(errList).To(BeNil())
			Ω(list).Should(ConsistOf("key1", "key2", "key3"))
		})
	})

	Describe("Test Delete() functionality", func() {
		It("should be possible to delete an existing value", func() {
			errInit := store.Init()
			Ω(errInit).Should(BeNil())

			errPut1 := store.Put(JobSessionType, "key1", "value1")
			Ω(errPut1).Should(BeNil())

			errPut2 := store.Put(JobSessionType, "key2", "value2")
			Ω(errPut2).Should(BeNil())

			exists := store.Exists(JobSessionType, "key1")
			Ω(exists).Should(BeTrue())

			errDel := store.Delete(JobSessionType, "key1")
			Ω(errDel).Should(BeNil())
			exists = store.Exists(JobSessionType, "key1")

			Ω(exists).Should(BeFalse())

			store.Exit()
		})

		It("should be error when deleting an non-existing value", func() {
			errInit := store.Init()

			Ω(errInit).Should(BeNil())
			errPut1 := store.Put(JobSessionType, "key1", "value1")
			Ω(errPut1).Should(BeNil())
			errPut2 := store.Put(JobSessionType, "key2", "value2")
			Ω(errPut2).Should(BeNil())

			errDel := store.Delete(JobSessionType, "key3")
			Ω(errDel).ShouldNot(BeNil())

			store.Exit()
		})
	})

	Describe("Testing Exists() functionality", func() {
		It("should signal the existence or non-existence of a key", func() {
			errInit := store.Init()
			Expect(errInit).To(BeNil())
			defer store.Exit()

			errPut1 := store.Put(JobSessionType, "key1", "value1")
			Expect(errPut1).To(BeNil())
			errPut2 := store.Put(JobSessionType, "key2", "value2")
			Expect(errPut2).To(BeNil())
			errPut3 := store.Put(JobSessionType, "key3", "")
			Expect(errPut3).To(BeNil())

			exists := store.Exists(JobSessionType, "key2")
			Expect(exists).Should(BeTrue())

			exists2 := store.Exists(ReservationSessionType, "key2")
			Expect(exists2).Should(BeFalse())

			exists3 := store.Exists(JobSessionType, "keyX")
			Expect(exists3).Should(BeFalse())

			exists4 := store.Exists(JobSessionType, "key3", "")
			Expect(exists4).To(BeTrue())
		})
	})

})
