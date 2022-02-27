package helper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helper Suite")
}
