package sidecar

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type JobContainerConfig struct {
	ClientSet     kubernetes.Interface
	Namespace     string
	PodName       string
	ContainerName string
}

// JobLifecycleSupervisor is a controller which watches a container
// in the same pod and does things when the container's state changes.
type JobLifecycleSupervisor struct {
	jc                 JobContainerConfig
	startupHook        map[int]func() error
	epilogHook         map[int]func(JobContainerConfig) error
	startupHooksResult map[int]error
	epilogHooksResult  map[int]error
}

// NewJobLifecylceSupervisor creates a new supervisor which executes scripts, then waits
// until the specified container in the given namespace and pod is finsihed, and then
// executes another set of specified functions.
func NewJobLifecylceSupervisor(jc JobContainerConfig) (*JobLifecycleSupervisor, error) {
	// if no client set is provided, create one with the assumption
	// to run as container inside a pod.
	if jc.ClientSet == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		jc.ClientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}
	return &JobLifecycleSupervisor{
		jc:                 jc,
		startupHook:        make(map[int]func() error),
		epilogHook:         make(map[int]func(JobContainerConfig) error),
		startupHooksResult: make(map[int]error),
		epilogHooksResult:  make(map[int]error),
	}, nil
}

// RegisterStartupHook registers a function which is executed when Run() is
// called. If RegisterStartupHook is called multiple times, the registered
// functions are called in the same order as they are registered.
func (l *JobLifecycleSupervisor) RegisterStartupHook(f func() error) error {
	l.startupHook[len(l.startupHook)] = f
	return nil
}

// RegisterEpilogHook registers a function which is executed the supervised
// container is finished. It is expected that the container does not re-run
// itself. If RegisterEpilogHook is called multiple times, the registered
// functions are called in the same order as they are registered.
func (l *JobLifecycleSupervisor) RegisterEpilogHook(f func(jc JobContainerConfig) error) error {
	l.epilogHook[len(l.epilogHook)] = f
	return nil
}

// Run starts the supervisor. First it runs all functioned registered as
// startup hooks. Then it waits until the supervised container is finished
// and then it runs all epilog hooks.
func (l *JobLifecycleSupervisor) Run() error {
	log.Println("Starting prolog scripts")
	// run all startup functions
	for i := 0; i < len(l.startupHook); i++ {
		l.startupHooksResult[i] = l.startupHook[i]()
	}

	log.Println("Waiting until batch job is finished")
	// wait until job is finished
	_, err := watchWaitTerminated(context.Background(), l.jc.ClientSet, l.jc.Namespace, l.jc.PodName, l.jc.ContainerName)
	if err != nil {
		return err
	}
	// job is finished

	log.Println("Running epilog scripts")
	// run all epilogs
	for i := 0; i < len(l.epilogHook); i++ {
		l.epilogHooksResult[i] = l.epilogHook[i](l.jc)
	}
	return nil
}

// watchWaitTerminated waits until the job container is finished and returns the execit code
func watchWaitTerminated(ctx context.Context, cs kubernetes.Interface, namespace, podName, containerName string) (int32, error) {
	pods := cs.CoreV1().Pods(namespace)
	watch, err := pods.Watch(ctx, v1.ListOptions{FieldSelector: "metadata.name=" + podName})
	if err != nil {
		log.Printf("watch: did not find container (could be finished already): %v\n", err)
		return 0, nil
	}
	defer watch.Stop()

	// check if pod is in an end state
	jobPod, err := pods.Get(ctx, podName, v1.GetOptions{})
	if err != nil {
		log.Printf("watch: failed getting pod %s: %v\n", podName, err)
	} else {
		log.Printf("watch: found pod\n")
		for _, status := range jobPod.Status.ContainerStatuses {
			fmt.Printf("checking container: %s\n", status.Name)
			if status.Name == containerName {
				log.Println("found container")
				if status.State.Terminated != nil {
					return status.State.Terminated.ExitCode, nil
				}
			}
		}
	}
	log.Printf("watch: waiting for pod events\n")
	for event := range watch.ResultChan() {
		log.Printf("watch: got pod change event: %v\n", event)
		pod, podType := event.Object.(*corev1.Pod)
		if !podType {
			log.Printf("internal error - watch got not a pod object")
			return -1, err
		}
		for _, status := range pod.Status.ContainerStatuses {
			if status.Name == containerName {
				if status.State.Terminated != nil {
					return status.State.Terminated.ExitCode, nil
				}
			}
		}
	}
	return -1, fmt.Errorf("internal error: unreachable code in watchWaitTerminated")
}
