package libdrmaa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLibdrmaa(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Libdrmaa Suite")
}
