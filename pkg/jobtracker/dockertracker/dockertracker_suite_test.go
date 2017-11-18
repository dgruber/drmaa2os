package dockertracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDockertracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockertracker Suite")
}
