package kubernetestracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
)

var _ = Describe("K8Connect", func() {

	Context("Clientset creation", func() {

		It("should be possible to create a new Clientset", func() {
			cs, err := NewClientSet()
			Ω(err).Should(BeNil())
			Ω(cs).ShouldNot(BeNil())
		})

		It("should create an error when .kube/config file is missing", func() {
			home := os.Getenv("HOME")
			os.Setenv("HOME", os.TempDir())
			defer os.Setenv("HOME", home)
			cs, err := NewClientSet()
			Ω(err).ShouldNot(BeNil())
			Ω(cs).Should(BeNil())
		})

	})

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
