package simpletracker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSimpletracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simpletracker Suite")
}
