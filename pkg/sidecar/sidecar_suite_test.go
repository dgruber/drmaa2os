package sidecar_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSidecar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sidecar Suite")
}
