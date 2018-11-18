package singularity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSingularity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Singularity Suite")
}
