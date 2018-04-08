package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

var _ = Describe("K8Connect", func() {

	Context("Helper functions", func() {

		It("Home directory should be returned and not empty", func() {
			home := homeDir()
			Ω(home).ShouldNot(Equal(""))
		})

		It("should create the path to the standard kubernetes config file", func() {
			cfg, err := kubeConfigFile()
			Ω(err).Should(BeNil())
			Ω(cfg).Should(ContainSubstring(".kube"))
			Ω(cfg).Should(ContainSubstring("config"))
		})

		It("should create a k8s client set (requires a kubernetes config)", func() {
			cs, err := NewClientSet()
			Ω(err).Should(BeNil())
			Ω(cs).ShouldNot(BeNil())
			cs, err = NewClientSet()
			Ω(err).Should(BeNil())
			Ω(cs).ShouldNot(BeNil())
		})

		Context("errors of helper functions", func() {

			It("should error when home path is empty", func() {
				originalHome := os.Getenv("HOME")
				originalUserprofile := os.Getenv("USERPROFILE")
				os.Setenv("HOME", "")
				os.Setenv("USERPROFILE", "")
				home := homeDir()
				Ω(home).Should(Equal(""))
				cfgFile, err := kubeConfigFile()
				Ω(err).ShouldNot(BeNil())
				Ω(cfgFile).Should(Equal(""))
				os.Setenv("HOME", originalHome)
				os.Setenv("USERPROFILE", originalUserprofile)
			})

		})

	})

})
