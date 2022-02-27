package d2hlp_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestD2hlp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "D2hlp Suite")
}
