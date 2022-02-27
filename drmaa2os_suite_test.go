package drmaa2os_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDrmaa2os(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Drmaa2os Suite")
}
