package dockertracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/api/types"
)

var _ = Describe("Convert", func() {

	Context("JobTemplate missing fields checks", func() {
		It("should recognize when JobCategory (container image name) is missing", func() {
			jt := drmaa2interface.JobTemplate{RemoteCommand: "/bin/sleep"}
			err := checkJobTemplate(jt)
			Ω(err).ShouldNot(BeNil())
		})

		It("should recognize when RemoteCommand is missing", func() {
			jt := drmaa2interface.JobTemplate{JobCategory: "image/image"}
			err := checkJobTemplate(jt)
			Ω(err).ShouldNot(BeNil())
		})

		It("should not fail when RemoteCommand and JobCategory is set", func() {
			jt := drmaa2interface.JobTemplate{JobCategory: "image/image", RemoteCommand: "/bin/sleep"}
			err := checkJobTemplate(jt)
			Ω(err).Should(BeNil())
		})
	})

	Context("Array Job ID convert functions", func() {

		It("should generate and resolve an array job ID into job IDs", func() {
			guids := []string{"1", "2", "3"}

			id := guids2ArrayJobID(guids)
			guidsOut, err := arrayJobID2GUIDs(id)

			Ω(err).Should(BeNil())
			Ω(guidsOut).Should(BeEquivalentTo(guids))
		})

	})

	Context("JobTemplate convert functions", func() {
		jt := drmaa2interface.JobTemplate{
			RemoteCommand:     "/bin/sleep",
			Args:              []string{"123"},
			JobCategory:       "my/image",
			WorkingDirectory:  "/working/dir",
			CandidateMachines: []string{"hostname"},
			StageInFiles:      map[string]string{"outer": "inner"},
		}

		It("should convert the JobTemplate settings", func() {
			cc, err := jobTemplateToContainerConfig(jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.Cmd[0]).Should(Equal("/bin/sleep"))
			Ω(cc.Cmd[1]).Should(Equal("123"))
			Ω(cc.Image).Should(Equal("my/image"))
			Ω(cc.WorkingDir).Should(Equal("/working/dir"))
			Ω(cc.Hostname).Should(Equal("hostname"))

			hc, err := jobTemplateToHostConfig(jt)
			Ω(err).Should(BeNil())
			Ω(hc.Binds).ShouldNot(BeNil())
			Ω(hc.Binds[0]).Should(Equal("outer:inner"))
		})

		It("should convert the environment variables", func() {
			jt.JobEnvironment = map[string]string{
				"env1": "value1",
				"env2": "value2",
			}
			cc, err := jobTemplateToContainerConfig(jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.Env).ShouldNot(BeNil())
			Ω(len(cc.Env)).Should(BeNumerically("==", 2))
			Ω(cc.Env[0]).Should(Or(Equal("env1=value1"), Equal("env2=value2")))
			Ω(cc.Env[1]).Should(Or(Equal("env1=value1"), Equal("env2=value2")))
		})
	})

	Context("JobState converter", func() {
		killed := &types.ContainerState{
			OOMKilled:  true,
			Dead:       false,
			ExitCode:   0,
			Paused:     false,
			Running:    false,
			Restarting: false,
		}

		dead := &types.ContainerState{
			OOMKilled:  false,
			Dead:       true,
			ExitCode:   1,
			Paused:     false,
			Running:    false,
			Restarting: false,
		}

		exit0 := &types.ContainerState{
			Status:     "exited",
			OOMKilled:  false,
			Dead:       true,
			ExitCode:   0,
			Paused:     false,
			Running:    false,
			Restarting: false,
		}

		exit1 := &types.ContainerState{
			Status:     "exited",
			OOMKilled:  false,
			Dead:       true,
			ExitCode:   1,
			Paused:     false,
			Running:    false,
			Restarting: false,
		}

		paused := &types.ContainerState{
			OOMKilled:  false,
			Dead:       false,
			ExitCode:   0,
			Paused:     true,
			Running:    false,
			Restarting: false,
		}

		restarting := &types.ContainerState{
			OOMKilled:  false,
			Dead:       false,
			ExitCode:   0,
			Paused:     false,
			Running:    true,
			Restarting: true,
		}

		running := &types.ContainerState{
			OOMKilled:  false,
			Dead:       false,
			ExitCode:   0,
			Paused:     false,
			Running:    true,
			Restarting: false,
		}

		It("should convert the state according to the documentation", func() {
			Ω(containerToDRMAA2State(killed)).Should(Equal(drmaa2interface.Failed))
			Ω(containerToDRMAA2State(dead)).Should(Equal(drmaa2interface.Failed))
			Ω(containerToDRMAA2State(exit0)).Should(Equal(drmaa2interface.Done))
			Ω(containerToDRMAA2State(exit1)).Should(Equal(drmaa2interface.Failed))
			Ω(containerToDRMAA2State(paused)).Should(Equal(drmaa2interface.Suspended))
			Ω(containerToDRMAA2State(restarting)).Should(Equal(drmaa2interface.Queued))
			Ω(containerToDRMAA2State(running)).Should(Equal(drmaa2interface.Running))
		})

	})

})
