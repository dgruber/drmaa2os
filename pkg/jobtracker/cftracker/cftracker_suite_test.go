package cftracker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCftracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cftracker Suite")
}
