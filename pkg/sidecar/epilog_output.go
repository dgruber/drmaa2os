package sidecar

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NewJobOutputToConfigMapEpilog returns an epilog function which gets the
// job output of a finished container and stores it in a ConfigMap named
// _jobname_-output-configmap.
// The stored output of the job is limited to 16MB.
func NewJobOutputToConfigMapEpilog() func(JobContainerConfig) error {
	return func(jc JobContainerConfig) error {
		limit := int64(16 * 1024 * 1024)
		if jc.ClientSet == nil {
			return errors.New("cannot get stdout of job because ClientSet is nil")
		}
		req := jc.ClientSet.CoreV1().Pods(jc.Namespace).GetLogs(jc.PodName, &v1.PodLogOptions{
			Container:  jc.ContainerName,
			Previous:   true,
			LimitBytes: &limit,
		})
		if req == nil {
			return errors.New("could not get request for log stream")
		}
		rc, err := req.Stream(context.Background())
		if err != nil {
			return err
		}
		output, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}
		// fake seems to return "fake logs"
		fmt.Printf("Output of job: %v\n", string(output))
		rc.Close()
		return CreateConfigMap(jc.ClientSet, jc.Namespace, jc.ContainerName+"-output-configmap", output)
	}
}

// CreateConfigMap creates a new configmap with the given data set in the "output" field of
// the ConfigMap.
func CreateConfigMap(cs kubernetes.Interface, nameSpace, name string, output []byte) error {
	_, err := cs.CoreV1().ConfigMaps(nameSpace).Create(
		context.Background(),
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nameSpace,
			},
			Data: map[string]string{
				"output": string(output),
			},
		}, metav1.CreateOptions{})
	return err
}
