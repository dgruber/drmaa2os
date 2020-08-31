package dockertracker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"time"
)

var _ = Describe("Convert", func() {

	Context("Internal helper functions", func() {
		It("should return a new PortSet", func() {
			// ip:public:private/proto
			ps := newPortSet("8080/tcp,1301/tcp")
			Ω(ps).ShouldNot(BeNil())
			Ω(len(ps)).Should(BeNumerically("==", 2))

			ps = newPortSet("192.168.2.111:80:6444/tcp,192.168.2.11:8080:6445/tcp")
			Ω(ps).ShouldNot(BeNil())
			Ω(len(ps)).Should(BeNumerically("==", 2))

			ps = newPortSet("80:6444/tcp,8080:6445/tcp")
			Ω(ps).ShouldNot(BeNil())
			Ω(len(ps)).Should(BeNumerically("==", 2))

		})

		It("should return a new PortMap", func() {
			pm := newPortBindings("8080/tcp,1301/tcp")
			Ω(pm).ShouldNot(BeNil())
			Ω(len(pm)).Should(BeNumerically("==", 2))

			pm = newPortBindings("192.168.2.111:80:6444/tcp,192.168.2.11:8080:6445/tcp")
			Ω(pm).ShouldNot(BeNil())
			Ω(len(pm)).Should(BeNumerically("==", 2))

			pm = newPortBindings("80:6444/tcp,8080:6445/tcp")
			Ω(pm).ShouldNot(BeNil())
			Ω(len(pm)).Should(BeNumerically("==", 2))
		})

		It("should return nil with wrong syntax for exposedPorts", func() {
			pm := newPortBindings("808-0/tcp,13+01/tcp")
			Ω(pm).Should(BeNil())
			ps := newPortSet("808-0/tcp,13+01/tcp")
			Ω(ps).Should(BeNil())
		})

		It("should return true when job state is in any given state", func() {
			match := isInExpectedState(drmaa2interface.Done, drmaa2interface.Failed, drmaa2interface.Done)
			Ω(match).Should(BeTrue())
			match = isInExpectedState(drmaa2interface.Done, drmaa2interface.Done)
			Ω(match).Should(BeTrue())
			match = isInExpectedState(drmaa2interface.Done, drmaa2interface.Done, drmaa2interface.Undetermined)
			Ω(match).Should(BeTrue())
		})

		It("should return false when job state is not in any given state", func() {
			match := isInExpectedState(drmaa2interface.Done, drmaa2interface.Failed, drmaa2interface.Undetermined)
			Ω(match).Should(BeFalse())
			match = isInExpectedState(drmaa2interface.Failed, drmaa2interface.Done)
			Ω(match).Should(BeFalse())
			match = isInExpectedState(drmaa2interface.Undetermined, drmaa2interface.Done, drmaa2interface.Failed)
			Ω(match).Should(BeFalse())
		})

	})

	Context("JobTemplate missing fields checks", func() {
		It("should recognize when JobCategory (container image name) is missing", func() {
			jt := drmaa2interface.JobTemplate{RemoteCommand: "/bin/sleep"}
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
			cc, err := jobTemplateToContainerConfig("session", jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.Cmd[0]).Should(Equal("/bin/sleep"))
			Ω(cc.Cmd[1]).Should(Equal("123"))
			Ω(cc.Image).Should(Equal("my/image"))
			Ω(cc.WorkingDir).Should(Equal("/working/dir"))
			Ω(cc.Hostname).Should(Equal("hostname"))
			jsLabel, exists := cc.Labels["drmaa2_jobsession"]
			Ω(exists).Should(BeTrue())
			Ω(jsLabel).Should(Equal("session"))
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
			cc, err := jobTemplateToContainerConfig("jobsession", jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.Env).ShouldNot(BeNil())
			Ω(len(cc.Env)).Should(BeNumerically("==", 2))
			Ω(cc.Env[0]).Should(Or(Equal("env1=value1"), Equal("env2=value2")))
			Ω(cc.Env[1]).Should(Or(Equal("env1=value1"), Equal("env2=value2")))
		})

		It("should set the user extension of JobTemplate", func() {
			jt.ExtensionList = map[string]string{"user": "testuser"}
			cc, err := jobTemplateToContainerConfig("", jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.User).Should(Equal("testuser"))
			Ω(cc.ExposedPorts).Should(BeNil())
		})

		It("should set the exposed ports when set in JobTemplate", func() {
			jt.ExtensionList = map[string]string{"exposedPorts": "80:6445/tcp"}
			cc, err := jobTemplateToContainerConfig("", jt)
			Ω(err).Should(BeNil())
			Ω(cc).ShouldNot(BeNil())
			Ω(cc.ExposedPorts).ShouldNot(BeNil())
			portSet := cc.ExposedPorts
			Ω(len(portSet)).Should(BeNumerically("==", 1))
			Ω(portSet).Should(HaveKey(nat.Port("6445/tcp")))
		})

		It("should set the portBindings when exposedPorts set in JobTemplate", func() {
			jt.ExtensionList = map[string]string{"exposedPorts": "80:6445/tcp"}
			jt.ExtensionList["restart"] = "unless-stopped"
			jt.ExtensionList["net"] = "host"
			jt.ExtensionList["privileged"] = "true"
			jt.ExtensionList["ipc"] = "host"
			jt.ExtensionList["uts"] = "host"
			jt.ExtensionList["pid"] = "host"
			jt.ExtensionList["rm"] = "true"
			hc, err := jobTemplateToHostConfig(jt)
			Ω(err).Should(BeNil())
			Ω(hc).ShouldNot(BeNil())
			Ω(len(hc.PortBindings)).Should(BeNumerically("==", 1))
			Ω(hc.PortBindings).Should(HaveKey(nat.Port("6445/tcp")))
			Ω(hc.RestartPolicy.IsUnlessStopped()).Should(BeTrue())
			Ω(hc.NetworkMode.IsHost()).Should(BeTrue())
			Ω(hc.Privileged).Should(BeTrue())
			Ω(hc.AutoRemove).Should(BeTrue())
			Ω(hc.IpcMode.IsHost()).Should(BeTrue())
			Ω(hc.UTSMode.IsHost()).Should(BeTrue())
			Ω(hc.PidMode.IsHost()).Should(BeTrue())
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

	Context("Container to JobInfo", func() {

		It("should convert basic information", func() {
			c := types.ContainerJSON{}
			c.ContainerJSONBase = &types.ContainerJSONBase{}
			c.ID = "jobid"
			c.Config = &container.Config{
				Hostname: "hostname",
			}

			finished, err := time.Parse(time.RFC3339Nano, "2015-01-06T15:47:32.080254511Z")
			Ω(err).Should(BeNil())
			started, err := time.Parse(time.RFC3339Nano, "2015-01-06T15:47:32.072697474Z")
			Ω(err).Should(BeNil())
			created, err := time.Parse(time.RFC3339Nano, "2015-01-06T15:47:31.485331387Z")
			Ω(err).Should(BeNil())

			c.State = &types.ContainerState{
				Status:     "exited",
				ExitCode:   13,
				FinishedAt: finished.Format(time.RFC3339Nano),
				StartedAt:  started.Format(time.RFC3339Nano),
			}
			c.Created = created.Format(time.RFC3339Nano)

			ji, _ := containerToDRMAA2JobInfo(c)
			Ω(ji.ID).Should(Equal("jobid"))
			Ω(ji.AllocatedMachines[0]).Should(Equal("hostname"))
			Ω(ji.State).Should(BeNumerically("==", drmaa2interface.Failed))
			Ω(ji.ExitStatus).Should(BeNumerically("==", 13))
			Ω(ji.SubmissionTime).Should(BeTemporally("==", created))
			Ω(ji.DispatchTime).Should(BeTemporally("==", started))
			Ω(ji.FinishTime).Should(BeTemporally("==", finished))
		})

	})

})
