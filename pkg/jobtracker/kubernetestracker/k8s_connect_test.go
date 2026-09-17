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
			if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
				Skip("the in-cluster configuration is used instead of .kube/config")
			}
			// an empty KUBECONFIG is ignored like an unset one
			GinkgoT().Setenv("KUBECONFIG", "")
			GinkgoT().Setenv("USERPROFILE", "")
			GinkgoT().Setenv("HOME", GinkgoT().TempDir())
			cs, err := NewClientSet()
			Expect(err).ShouldNot(BeNil())
			Expect(cs).Should(BeNil())
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

			BeforeEach(func() {
				// GinkgoT().Setenv restores the environment after each spec
				GinkgoT().Setenv("USERPROFILE", "")
				GinkgoT().Setenv("KUBECONFIG", "")
				GinkgoT().Setenv("HOME", GinkgoT().TempDir())
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
