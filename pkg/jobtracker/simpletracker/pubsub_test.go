package simpletracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Pubsub", func() {

	Context("Single producer and single consumer", func() {

		It("should return the running state of the job", func() {
			// producer
			ps, jeCh := NewPubSub()

			ps.StartBookKeeper()

			// consumer
			waitCh, err := ps.Register("13", drmaa2interface.Running)

			Ω(err).To(BeNil())

			// produce
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "13"}

			// consume
			evt := <-waitCh

			Ω(evt).To(Equal(drmaa2interface.Running))
		})

	})

	Context("Multiple producer and single consumer", func() {

		It("should detect all running states from the jobs", func() {

			// producer
			ps, jeCh := NewPubSub()

			ps.StartBookKeeper()

			// consumer
			waitCh, err := ps.Register("13", drmaa2interface.Running)

			Ω(err).To(BeNil())

			// produce
			jeCh <- JobEvent{JobState: drmaa2interface.Suspended, JobID: "11"}
			jeCh <- JobEvent{JobState: drmaa2interface.Requeued, JobID: "12"}
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "13"}
			jeCh <- JobEvent{JobState: drmaa2interface.QueuedHeld, JobID: "14"}

			// consume
			//evt := <-waitCh

			//Ω(evt).To(Equal(drmaa2interface.Running))
			Eventually(waitCh).Should(Receive(Equal(drmaa2interface.Running)))

		})
	})

	Context("Single producer and multiple consumer", func() {

		It("should send to all consumers the running state of the job", func() {

			// producer
			ps, jeCh := NewPubSub()

			ps.StartBookKeeper()

			// consumer
			waitCh, err := ps.Register("13", drmaa2interface.Running)
			Ω(err).Should(BeNil())

			// consumer
			waitCh2, err2 := ps.Register("13", drmaa2interface.Running)
			Ω(err2).Should(BeNil())

			// consumer
			waitCh3, err3 := ps.Register("13", drmaa2interface.Running)
			Ω(err3).Should(BeNil())

			// consumer
			waitCh4, err4 := ps.Register("13", drmaa2interface.Queued)
			Ω(err4).Should(BeNil())

			// produce
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "13"}

			// consume
			Eventually(waitCh3).Should(Receive(Equal(drmaa2interface.Running)))
			Eventually(waitCh2).Should(Receive(Equal(drmaa2interface.Running)))
			Eventually(waitCh).Should(Receive(Equal(drmaa2interface.Running)))

			// this channel should time out
			Consistently(waitCh4).ShouldNot(Receive())

		})
	})

	Context("Multiple producer and multiple consumer", func() {

		It("should send to all consumers the running state of the job", func() {

			// producer
			ps, jeCh := NewPubSub()

			ps.StartBookKeeper()

			// consumer
			waitCh, err := ps.Register("14", drmaa2interface.Running)
			Ω(err).Should(BeNil())

			// consumer
			waitCh2, err2 := ps.Register("13", drmaa2interface.Running)
			Ω(err2).Should(BeNil())

			// consumer
			waitCh3, err3 := ps.Register("13", drmaa2interface.Failed)
			Ω(err3).Should(BeNil())

			// consumer
			waitCh4, err4 := ps.Register("13", drmaa2interface.Requeued)
			Ω(err4).Should(BeNil())

			// consumer
			waitCh5, err5 := ps.Register("1", drmaa2interface.Requeued)
			Ω(err5).Should(BeNil())

			// produce
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "13"}
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "14"}
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "15"}
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "16"}
			jeCh <- JobEvent{JobState: drmaa2interface.Running, JobID: "17"}
			jeCh <- JobEvent{JobState: drmaa2interface.Requeued, JobID: "13"}
			jeCh <- JobEvent{JobState: drmaa2interface.Failed, JobID: "13"}

			// consume
			Eventually(waitCh4).Should(Receive(Equal(drmaa2interface.Requeued)))
			Eventually(waitCh3).Should(Receive(Equal(drmaa2interface.Failed)))
			Eventually(waitCh2).Should(Receive(Equal(drmaa2interface.Running)))
			Eventually(waitCh).Should(Receive(Equal(drmaa2interface.Running)))

			// this channel should time out
			Consistently(waitCh5).ShouldNot(Receive())

		})
	})

})
