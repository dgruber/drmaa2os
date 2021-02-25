package main

import (
	"log"
	"os"
	"strings"

	"github.com/dgruber/drmaa2os/pkg/sidecar"
)

// Sidecar which can run along with the Kubernetes job to provide
// additional functionality which can't be implemented by other
// means.
func main() {
	// requires downward API in the pod definition set to
	// expose podname and namespace
	podName := os.Getenv("DRMAA2OS_POD_NAME")
	log.Printf("Pod name: %s\n", podName)

	podNamespace := os.Getenv("DRMAA2OS_POD_NAMESPACE")
	log.Printf("Pod namespace: %s\n", podNamespace)

	// In DRMAA2 jobs the job name is the pod name + the postfix
	// for the pod added by the Kubernetes job controller. The
	// container name is the same as the job name. So we get the
	// job container name out of the pod name removing the postfix
	// appended with -XYZ.
	ctrStr := strings.Split(podName, "-")
	if len(ctrStr) < 2 {
		log.Fatal("podname does not contain -")
	}
	containerName := strings.Join(ctrStr[:len(ctrStr)-1], "-")
	log.Printf("Container ID: %s\n", containerName)

	lm, err := sidecar.NewJobLifecylceSupervisor(
		sidecar.JobContainerConfig{
			ClientSet:     nil,
			Namespace:     podNamespace,
			PodName:       podName,
			ContainerName: containerName})

	if err != nil {
		log.Fatal(err)
	}
	lm.RegisterStartupHook(func() error {
		log.Printf("drmaa2os sidecar: startup")
		return nil
	})

	lm.RegisterEpilogHook(sidecar.NewJobOutputToConfigMapEpilog())
	err = lm.Run()
	if err != nil {
		log.Fatal(err)
	}
}
