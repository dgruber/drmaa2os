package libdrmaa

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCdrmaajobtracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cdrmaajobtracker Suite")
}
