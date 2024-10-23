package kubernetestracker

import (
	. "github.com/onsi/ginkgo/v2"
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
			// test only works when KUBECONFIG is unset (as this is the fallback)
			if _, exists := os.LookupEnv("KUBECONFIG"); exists {
				Ω(err).To(BeNil())
			} else {
				Ω(err).ShouldNot(BeNil())
				Ω(cs).Should(BeNil())
			}
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
			// the file should exist
			Ω(cfg).Should(BeAnExistingFile())
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

			var originalHome, originalUserprofile, kubeconfig string
			var dir string

			BeforeEach(func() {
				originalHome = os.Getenv("HOME")
				originalUserprofile = os.Getenv("USERPROFILE")
				kubeconfig = os.Getenv("KUBECONFIG")
				os.Setenv("USERPROFILE", "")
				os.Setenv("KUBECONFIG", "")
				var err error
				dir, err = os.MkdirTemp("", "test")
				Expect(err).Should(BeNil())
				os.Setenv("HOME", dir)
			})

			AfterEach(func() {
				os.Setenv("HOME", originalHome)
				os.Setenv("USERPROFILE", originalUserprofile)
				os.Setenv("KUBECONFIG", kubeconfig)
				if dir != "" {
					os.RemoveAll(dir)
				}
			})

			It("should error when home path is empty", func() {
				home := homeDir()
				Ω(home).ShouldNot(Equal(""))
				Expect(home).To(Equal(os.Getenv("HOME")))
				cfgFile, err := kubeConfigFile()
				Ω(err).ShouldNot(BeNil())
				Ω(cfgFile).Should(Equal(""))
			})

		})

	})

})
