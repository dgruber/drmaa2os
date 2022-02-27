package sidecar_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/dgruber/drmaa2os/pkg/sidecar"
	"k8s.io/client-go/kubernetes/fake"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

var _ = Describe("Lifecyclemanager", func() {

	Context("Basic functions", func() {

		It("should be possible to create a new LifecycleManager", func() {
			cs := fake.NewSimpleClientset()
			lm, err := NewJobLifecylceSupervisor(
				JobContainerConfig{
					ClientSet:     cs,
					Namespace:     "default",
					PodName:       "podname",
					ContainerName: "containerid"})
			Expect(err).To(BeNil())
			Expect(lm).NotTo(BeNil())
		})

		It("should run startup scripts and epilog scripts even when main batch job container is finished already", func() {
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
										FinishedAt: v1.NewTime(time.Now()),
									},
								},
								ContainerID: "docker://sidecarid123",
							},
						},
					},
				},
			)
			lm, err := NewJobLifecylceSupervisor(
				JobContainerConfig{
					ClientSet:     cs,
					Namespace:     "default",
					PodName:       "podname",
					ContainerName: "sidecar"})
			Expect(err).To(BeNil())
			Expect(lm).NotTo(BeNil())
			var startupCalled int
			var epilogCalled int
			err = lm.RegisterStartupHook(func() error {
				startupCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterStartupHook(func() error {
				startupCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
				epilogCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
				epilogCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.Run()
			Expect(err).To(BeNil())

			// check if both hooks are called
			Expect(startupCalled).To(BeNumerically("==", 2))
			Expect(epilogCalled).To(BeNumerically("==", 2))
		})

		It("should run startup scripts and epilog scripts even when main batch job container is finished already", func() {
			cs := fake.NewSimpleClientset(
				&corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "podname",
						Namespace: "default"},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "jobname",
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
								Name: "jobname",
								State: corev1.ContainerState{
									Terminated: &corev1.ContainerStateTerminated{
										FinishedAt: v1.NewTime(time.Now()),
										ExitCode:   77,
									},
								},
								ContainerID: "docker://sidecarid123",
							},
						},
					},
				},
			)
			lm, err := NewJobLifecylceSupervisor(
				JobContainerConfig{
					ClientSet:     cs,
					Namespace:     "default",
					PodName:       "podname",
					ContainerName: "jobname"})
			Expect(err).To(BeNil())
			Expect(lm).NotTo(BeNil())
			var startupCalled int
			var epilogCalled int
			err = lm.RegisterStartupHook(func() error {
				startupCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterStartupHook(func() error {
				startupCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
				epilogCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
				epilogCalled++
				return nil
			})
			Expect(err).To(BeNil())
			err = lm.Run()
			Expect(err).To(BeNil())

			// check if both hooks are called
			Expect(startupCalled).To(BeNumerically("==", 2))
			Expect(epilogCalled).To(BeNumerically("==", 2))
		})

	})

	It("should run s epilog scripts even when main batch job container is finished", func() {
		cs := fake.NewSimpleClientset(&corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name:      "podname",
				Namespace: "default"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "jobname",
						Image: "image",
					},
					{
						Name:  "sidecar",
						Image: "image",
					}},
			}})
		// no pods yet

		watcher := watch.NewFake()
		cs.PrependWatchReactor("pods", testing.DefaultWatchReactor(watcher, nil))

		lm, err := NewJobLifecylceSupervisor(
			JobContainerConfig{
				ClientSet:     cs,
				Namespace:     "default",
				PodName:       "podname",
				ContainerName: "jobname"})
		Expect(err).To(BeNil())
		Expect(lm).NotTo(BeNil())
		var epilogCalled int
		Expect(err).To(BeNil())
		err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
			epilogCalled++
			return nil
		})
		Expect(err).To(BeNil())
		err = lm.RegisterEpilogHook(func(jc JobContainerConfig) error {
			epilogCalled++
			return nil
		})
		Expect(err).To(BeNil())

		start := time.Now()
		waitUntilEvent := time.Millisecond * 500
		// add job event
		go func() {
			fmt.Printf("add job event\n")
			<-time.Tick(waitUntilEvent)
			watcher.Modify(&corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "podname",
					Namespace: "default"},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "jobname",
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
							Name: "jobname",
							State: corev1.ContainerState{
								Terminated: &corev1.ContainerStateTerminated{
									FinishedAt: v1.NewTime(time.Now()),
									ExitCode:   77,
								},
							},
							ContainerID: "docker://jobname123",
						},
					},
				},
			})
			fmt.Printf("changed object\n")
		}()

		err = lm.Run()
		Expect(err).To(BeNil())

		Expect(time.Now()).To(BeTemporally(">=", start.Add(waitUntilEvent)))
		Expect(epilogCalled).To(BeNumerically("==", 2))
	})

	Context("Error cases", func() {

	})

})
