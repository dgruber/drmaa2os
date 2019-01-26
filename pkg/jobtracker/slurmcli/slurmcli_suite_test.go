package slurmcli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSlurmcli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slurmcli Suite")
}
