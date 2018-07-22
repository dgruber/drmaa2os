package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("K8sState", func() {

	Context("Job state", func() {

		It("should return undetermined as state when job is not found", func() {
			cs, err := NewClientSet()
			Ω(err).Should(BeNil())
			state := DRMAA2State(cs.BatchV1().Jobs("default"), "doesnotexist")
			Ω(state).Should(Equal(drmaa2interface.Undetermined))
		})

	})

	Context("JobStatus conversion", func() {
		var status batchv1.JobStatus

		BeforeEach(func() {
			status = batchv1.JobStatus{
				Active:    0,
				Succeeded: 0,
				Failed:    0,
			}
		})

		It("should convert nil to Undetermined state", func() {
			Ω(convertJobStatus2JobState(nil)).Should(Equal(drmaa2interface.Undetermined))
		})

		It("should convert active to Running state", func() {
			status.Active = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Running))
		})

		It("should convert failed to Failed state", func() {
			status.Failed = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Failed))
		})

		It("should convert into failed state when it is not succeeded/failed/active but completed already", func() {
			completed := metav1.NewTime(time.Now().Add(-time.Second))
			status.CompletionTime = &completed
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Failed))
		})

		It("should convert succeeded to Done state", func() {
			status.Succeeded = 1
			Ω(convertJobStatus2JobState(&status)).Should(Equal(drmaa2interface.Done))
		})

		It("should convert unset states to Undetermined state", func() {
			var s batchv1.JobStatus
			Ω(convertJobStatus2JobState(&s)).Should(Equal(drmaa2interface.Undetermined))
		})
	})

})
