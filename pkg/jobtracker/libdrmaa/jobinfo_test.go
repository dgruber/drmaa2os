package libdrmaa

import (
	"fmt"
	"time"

	"github.com/dgruber/drmaa2interface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobinfo", func() {

	Context("basic tests", func() {

		It("should not crash when job info is not allocated", func() {
			info := ConvertDRMAAJobInfoToDRMAA2JobInfo(nil)
			Expect(info.ID).To(Equal(""))
		})

	})

	Context("Tests with job executions", func() {

		It("should set a meaningful submission, dispatch, and finish time", func() {
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			before := time.Now().Add(-time.Second)
			jobID, err := d.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobID).NotTo(Equal(""))

			// could interfere with other tests
			err = d.Wait(jobID, time.Second*30, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			ji, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())

			Expect(ji.SubmissionTime).To(BeTemporally(">", before))
			Expect(ji.SubmissionTime).To(BeTemporally("<", time.Now()))
			Expect(ji.DispatchTime).To(BeTemporally(">", before))
			Expect(ji.DispatchTime).To(BeTemporally("<", time.Now()))
			Expect(ji.FinishTime).To(BeTemporally(">", before))
			Expect(ji.FinishTime).To(BeTemporally("<", time.Now()))

			Expect(ji.SubmissionTime).To(BeTemporally("<=", ji.DispatchTime))
			Expect(ji.DispatchTime).To(BeTemporally("<=", ji.FinishTime))
		})

	})

	Context("Time", func() {

		It("should convert seconds string to time", func() {
			now := time.Now().Round(time.Second)
			secondsSinceEpoch := now.Unix()
			seconds := fmt.Sprintf("%d.0000", secondsSinceEpoch)
			convertedTime := ConvertUnixToTime(seconds)
			Expect(convertedTime).To(BeTemporally("==", now))
		})

	})

})
