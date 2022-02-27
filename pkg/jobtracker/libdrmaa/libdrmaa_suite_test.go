package libdrmaa

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLibdrmaa(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Libdrmaa Suite")
}
