package containerdtracker_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainerdtracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containerdtracker Suite")
}
