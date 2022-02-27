package podmantracker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
)

var _ = Describe("Run", func() {

	var jt drmaa2interface.JobTemplate

	BeforeEach(func() {
		jt = drmaa2interface.JobTemplate{
			RemoteCommand: "sleep",
			JobCategory:   "busybox:latest",
		}
	})

	Context("Container spec settings", func() {

		It("should allocate a PTY", func() {
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.Terminal).To(BeTrue())
		})

		It("should set command and args", func() {
			jt.Args = []string{"1", "2", "3"}
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.Command).To(Equal([]string{"sleep", "1", "2", "3"}))
		})

		It("should set environment variables", func() {
			jt.JobEnvironment = map[string]string{
				"first": "variable",
				"and":   "second",
			}
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.Env).To(ContainElement("variable"))
			Expect(spec.Env).To(ContainElement("second"))
		})

		It("should set the hostname according to the candidate machine", func() {
			jt.CandidateMachines = []string{"myhostname"}
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.Hostname).To(Equal("myhostname"))
		})

		It("should set the working directory inside the container", func() {
			jt.WorkingDirectory = "/application"
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.WorkDir).To(Equal("/application"))
		})

		It("should set all extensions", func() {
			jt.ExtensionList = map[string]string{
				"user":         "hans",
				"exposedPorts": "127.0.0.1:8181:80,8282:443",
				"privileged":   "true",
				"restart":      "always",
				"ipc":          "host",
				"uts":          "host",
				"pid":          "host",
				"rm":           "TRUE",
			}
			spec, err := CreateContainerSpec(jt)
			Expect(err).To(BeNil())
			Expect(spec.User).To(Equal("hans"))
			Expect(spec.PortMappings).NotTo(BeNil())
			Expect(len(spec.PortMappings)).To(BeNumerically("==", 2))
			Expect(spec.Privileged).To(BeTrue())
			Expect(spec.RestartPolicy).To(Equal("always"))
			Expect(spec.IpcNS.IsHost()).To(BeTrue())
			Expect(spec.UtsNS.IsHost()).To(BeTrue())
			Expect(spec.PidNS.IsHost()).To(BeTrue())
			Expect(spec.Remove).To(BeTrue())
		})

	})

})
