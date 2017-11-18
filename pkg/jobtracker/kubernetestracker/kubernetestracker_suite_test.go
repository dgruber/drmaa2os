package kubernetestracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKubernetestracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubernetestracker Suite")
}
