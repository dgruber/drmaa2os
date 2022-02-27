package server_test

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {

	Context("job state mapping", func() {

		It("should map job states between DRMAA2 and OpenAPI spec back and forth", func() {

			for _, state := range []drmaa2interface.JobState{
				drmaa2interface.Running, drmaa2interface.Done, drmaa2interface.Failed,
				drmaa2interface.Queued, drmaa2interface.QueuedHeld,
				drmaa2interface.Requeued, drmaa2interface.RequeuedHeld,
				drmaa2interface.Suspended,
				drmaa2interface.Unset,
				drmaa2interface.Undetermined} {
				Expect(client.ConvertJobStateToDRMAA2(string(ConvertJobState(state.String())))).To(Equal(state))
			}

		})

	})

})
