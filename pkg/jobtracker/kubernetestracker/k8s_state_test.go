package kubernetestracker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/dgruber/drmaa2interface"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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

	Context("JobStatus conversion", func() {

		It("should return the output of a finsihed job", func() {
			cs := fake.NewSimpleClientset(
				&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"job-name": "job",
						},
						Name:      "podname",
						Namespace: "default"},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "job",
								Image: "image",
							},
						},
					},
					Status: corev1.PodStatus{
						ContainerStatuses: []corev1.ContainerStatus{
							{
								Name: "job",
								State: corev1.ContainerState{
									Terminated: &corev1.ContainerStateTerminated{
										FinishedAt: metav1.NewTime(time.Now()),
									},
								},
								ContainerID: "docker://sidecarid123",
							},
						},
					},
				},
			)
			podList, err := GetPodsForJob(cs, "default", "job")
			Ω(err).Should(BeNil())
			podName := GetLastStartedPod(podList).Name
			output, err := GetJobOutput(cs, "default", "job", podName)
			Ω(err).Should(BeNil())
			Ω(string(output)).Should(Equal("fake logs"))
		})

	})

})
