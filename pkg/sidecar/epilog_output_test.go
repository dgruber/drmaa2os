package sidecar_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/dgruber/drmaa2os/pkg/sidecar"
)

var _ = Describe("EpilogOutput", func() {

	Context("Basic tests", func() {

		It("should return the output of a finished container", func() {

			cs := fake.NewSimpleClientset(
				&corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "podname",
						Namespace: "default"},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "name",
								Image: "image",
							},
							{
								Name:  "sidecar",
								Image: "image",
							}},
					},
					Status: corev1.PodStatus{
						ContainerStatuses: []corev1.ContainerStatus{
							{
								Name: "sidecar",
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

			f := NewJobOutputToConfigMapEpilog()
			err := f(JobContainerConfig{
				ClientSet:     cs,
				Namespace:     "default",
				PodName:       "podname",
				ContainerName: "sidecar",
			})
			Expect(err).To(BeNil())

		})

		It("should create a new ConfigMap", func() {
			fakeClientSet := fake.NewSimpleClientset()
			CreateConfigMap(fakeClientSet, "test", "joboutput", []byte("test data"))
			cm, err := fakeClientSet.CoreV1().ConfigMaps("test").Get(context.Background(), "joboutput", metav1.GetOptions{})
			Expect(err).To(BeNil())
			output, exists := cm.Data["output"]
			Expect(exists).To(BeTrue())
			Expect(output).To(Equal("test data"))
		})
	})

})
