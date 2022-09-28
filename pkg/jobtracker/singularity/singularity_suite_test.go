package singularity_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSingularity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Singularity Suite")
}
