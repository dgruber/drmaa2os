package kubernetestracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
	"testing"
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
	_, err := kubernetestracker.CreateClientSet()
	if err != nil {
		k8sChecked = true
		k8savailable = true
	} else {
		k8sChecked = true
		k8savailable = false
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
