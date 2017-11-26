package simpletracker_test

import (
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"
	"github.com/scalingdata/gosigar"

	"io/ioutil"
	"os"
	"strings"
	_ "time"
)

func exists(pid int) bool {
	state := sigar.ProcState{}
	err := state.Get(pid)
	if err != nil {
		return false
	}
	return true
}

var _ = Describe("OS specific functionality", func() {

	var (
		jt    drmaa2interface.JobTemplate
		outCh chan JobEvent
	)

	BeforeEach(func() {
		jt = drmaa2interface.JobTemplate{RemoteCommand: "sleep", Args: []string{"1"}}
		outCh = make(chan JobEvent, 1)
	})

	Context("Process state change when the process is in the expected state", func() {

		It("should be able to create a process", func() {
			pid, err := StartProcess("1", jt, outCh)

			Ω(err).Should(BeNil())
			Ω(pid).ShouldNot(BeNumerically("<=", 1))

			proc, errFind := os.FindProcess(pid)
			Ω(errFind).Should(BeNil())
			Ω(proc).ShouldNot(BeNil())
		})

		It("should be able to terminate a process", func() {
			pid, err := StartProcess("1", jt, outCh)

			Ω(err).Should(BeNil())
			Ω(pid).ShouldNot(BeNumerically("<=", 1))

			proc, errFind := os.FindProcess(pid)
			Ω(errFind).Should(BeNil())
			Ω(proc).ShouldNot(BeNil())

			errKill := KillPid(pid)
			Ω(errKill).Should(BeNil())
		})

		It("should be possible to suspend and resume a process", func() {
			pid, err := StartProcess("1", jt, outCh)

			Ω(err).Should(BeNil())

			err = SuspendPid(pid)
			Ω(err).Should(BeNil())

			err = ResumePid(pid)
			Ω(err).Should(BeNil())
		})
	})

	Context("Redirection of file descriptors", func() {

		It("should be possible to redirect stdout to file", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "/bin/echo"
			jt.Args = []string{"output"}
			jt.OutputPath = file.Name()
			_, err = StartProcess("1", jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh

			out, err := ioutil.ReadFile(file.Name())
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("output"))

			os.Remove(file.Name())
		})

		It("should be possible to redirect stdin from file", func() {
			fileOut, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())
			fileIn, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			_, err = fileIn.WriteString("inout\n")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "./stdin.sh"

			jt.InputPath = fileIn.Name()
			jt.OutputPath = fileOut.Name()

			_, err = StartProcess("1", jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh

			out, err := ioutil.ReadFile(fileOut.Name())
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("inout"))

			os.Remove(fileIn.Name())
			os.Remove(fileOut.Name())
		})

		It("should be possible to redirect stderr to file", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "./stderr.sh"

			jt.OutputPath = file.Name()
			_, err = StartProcess("1", jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh

			out, err := ioutil.ReadFile(file.Name())
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("error"))

			os.Remove(file.Name())
		})

	})

	Context("Redirection of file descriptors sad path", func() {

		It("should not be possible to point stdin and stdout / stderr to same file", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "/bin/echo"
			jt.Args = []string{"output"}

			jt.InputPath = file.Name()
			jt.OutputPath = file.Name()
			_, err = StartProcess("1", jt, outCh)

			Ω(err).ShouldNot(BeNil())

			jt.OutputPath = ""
			jt.ErrorPath = file.Name()
			_, err = StartProcess("1", jt, outCh)

			Ω(err).ShouldNot(BeNil())
		})
	})

})
