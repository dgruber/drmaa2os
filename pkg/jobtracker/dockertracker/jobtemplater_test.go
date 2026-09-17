package dockertracker_test

import (
	"encoding/base64"
	"encoding/json"

	"github.com/dgruber/drmaa2interface"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobtemplater without Docker daemon", func() {

	var daemon *fakeDockerDaemon
	var jobTemplate drmaa2interface.JobTemplate

	BeforeEach(func() {
		jobTemplate = drmaa2interface.JobTemplate{
			RemoteCommand: "/bin/sleep",
			Args:          []string{"1"},
			JobCategory:   "alpine",
		}
		encoded, err := json.Marshal(jobTemplate)
		Ω(err).Should(BeNil())

		labels := map[string]string{
			ContainerLabelJobTemplate: base64.StdEncoding.EncodeToString(encoded),
		}
		daemon = startFakeDockerDaemon("1.44", "1.44", labels)
		daemon.useDaemon("")
	})

	It("should return the job template of a container", func() {
		tracker, err := New("")
		Ω(err).Should(BeNil())

		jt, err := tracker.JobTemplate("container1")
		Ω(err).Should(BeNil())
		Ω(jt).Should(Equal(jobTemplate))
	})

	It("should reuse the connection of the tracker", func() {
		tracker, err := New("")
		Ω(err).Should(BeNil())

		for i := 0; i < 3; i++ {
			_, err = tracker.JobTemplate("container1")
			Ω(err).Should(BeNil())
		}
		Ω(daemon.pings.Load()).Should(BeNumerically("==", 1))
	})

	It("should return an error when the tracker is not initialized", func() {
		var tracker DockerTracker
		_, err := tracker.JobTemplate("container1")
		Ω(err).ShouldNot(BeNil())
	})

	It("should return an error when the label is missing", func() {
		startFakeDockerDaemon("1.44", "1.44", nil).useDaemon("")
		tracker, err := New("")
		Ω(err).Should(BeNil())

		_, err = tracker.JobTemplate("container1")
		Ω(err).ShouldNot(BeNil())
	})

	It("should read the job template without a tracker", func() {
		jt, err := ReadJobTemplateFromLabel("container1")
		Ω(err).Should(BeNil())
		Ω(jt).Should(Equal(jobTemplate))
	})

})
