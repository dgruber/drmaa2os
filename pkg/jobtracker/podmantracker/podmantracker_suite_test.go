package podmantracker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPodmantracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Podmantracker Suite")
}
