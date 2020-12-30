package kubernetestracker_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
)

var k8sChecked bool = false
var k8savailable bool = false

func TestKubernetestracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubernetestracker Suite")
}

func k8sIsAvailable() bool {
	if k8sChecked {
		return k8savailable
	}
	_, err := kubernetestracker.NewClientSet()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		k8sChecked = true
		k8savailable = false
	} else {
		k8sChecked = true
		k8savailable = true
	}
	return k8savailable
}

func WhenK8sIsAvailableIt(description string, f interface{}) {
	if k8sIsAvailable() {
		It(description, f)
	} else {
		PIt(description, f)
	}
}

func WhenK8sIsAvailableFIt(description string, f interface{}) {
	if k8sIsAvailable() {
		FIt(description, f)
	} else {
		PIt(description, f)
	}
}
