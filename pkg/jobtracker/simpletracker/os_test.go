package simpletracker_test

import (
	"fmt"

	sigar "github.com/cloudfoundry/gosigar"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"

	"io/ioutil"
	"os"
	"strings"
	_ "time"
)

func exists(pid int) bool {
	state := sigar.ProcState{}
	err := state.Get(pid)
	return err == nil
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
			pid, err := StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())
			Ω(pid).ShouldNot(BeNumerically("<=", 1))

			proc, errFind := os.FindProcess(pid)
			Ω(errFind).Should(BeNil())
			Ω(proc).ShouldNot(BeNil())
		})

		It("should be able to terminate a process", func() {
			pid, err := StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())
			Ω(pid).ShouldNot(BeNumerically("<=", 1))

			proc, errFind := os.FindProcess(pid)
			Ω(errFind).Should(BeNil())
			Ω(proc).ShouldNot(BeNil())

			errKill := KillPid(pid)
			Ω(errKill).Should(BeNil())
		})

		It("should be possible to suspend and resume a process", func() {
			pid, err := StartProcess("1", 0, jt, outCh)
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
			_, err = StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh
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

			_, err = StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh
			<-outCh

			out, err := ioutil.ReadFile(fileOut.Name())
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("inout"))

			os.Remove(fileIn.Name())
			os.Remove(fileOut.Name())
		})

		It("should be possible to redirect stdin from file (2)", func() {
			fileOut, err := ioutil.TempFile("", "d2ostest")
			Ω(err).Should(BeNil())
			fileOutName := fileOut.Name()
			fileOut.Close()

			fileIn, err := ioutil.TempFile("", "d2ostest")
			Ω(err).Should(BeNil())
			_, err = fileIn.WriteString("inout\n")
			Ω(err).Should(BeNil())
			fileInName := fileIn.Name()
			fileIn.Close()

			// sort expects that input is closed
			jt.RemoteCommand = "/bin/cat"
			jt.Args = nil

			jt.InputPath = fileInName
			jt.OutputPath = fileOutName

			_, err = StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh
			<-outCh

			out, err := ioutil.ReadFile(fileOutName)
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("inout"))

			os.Remove(fileInName)
			os.Remove(fileOutName)
		})

		It("should be possible to redirect stderr to file", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "./stderr.sh"
			jt.ErrorPath = file.Name()
			_, err = StartProcess("1", 0, jt, outCh)
			Ω(err).Should(BeNil())

			<-outCh
			<-outCh

			out, err := ioutil.ReadFile(file.Name())
			Ω(err).Should(BeNil())
			Ω(strings.TrimSpace(string(out))).Should(Equal("error"))

			os.Remove(file.Name())
		})

	})

	Context("Redirection of file descriptors errors", func() {

		It("should not be possible to point stdin and stdout / stderr to same file", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())

			jt.RemoteCommand = "/bin/echo"
			jt.Args = []string{"output"}

			jt.InputPath = file.Name()
			jt.OutputPath = file.Name()
			_, err = StartProcess("1", 0, jt, outCh)

			Ω(err).ShouldNot(BeNil())

			jt.OutputPath = ""
			jt.ErrorPath = file.Name()
			_, err = StartProcess("1", 0, jt, outCh)

			Ω(err).ShouldNot(BeNil())
		})

		It("should return an error when stderr file can not be generated", func() {
			file, err := ioutil.TempFile(os.TempDir(), "d2ostest")
			Ω(err).Should(BeNil())
			file.Close()

			jt.RemoteCommand = "/bin/echo"
			jt.Args = []string{"output"}

			jt.ErrorPath = "/non/existing/path"
			_, err = StartProcess("1", 0, jt, outCh)

			// wrong path / can not create file
			Ω(err).ShouldNot(BeNil())

			jt.InputPath = "non/existing/path"
			jt.ErrorPath = "/dev/stderr"

			_, err = StartProcess("1", 0, jt, outCh)

			// wrong path / can not create file
			Ω(err).ShouldNot(BeNil())

			jt.InputPath = ""
			jt.OutputPath = ""
			jt.ErrorPath = file.Name()
			_, err = StartProcess("1", 0, jt, outCh)

			Ω(err).Should(BeNil())
		})

	})

	Context("Potential race conditions", func() {

		It("should not block", func() {

			outCh = make(chan JobEvent, 3000)

			jt.RemoteCommand = "echo"
			jt.Args = []string{"x"}
			jt.OutputPath = "/dev/stdout"
			jt.ErrorPath = "/dev/stderr"

			for i := 0; i < 100; i++ {
				_, err := StartProcess(fmt.Sprintf("%d", i), 0, jt, outCh)
				Ω(err).Should(BeNil())
			}

			done := 0
			other := 0
			for jobEvent := range outCh {
				if jobEvent.JobInfo.State == drmaa2interface.Done {
					done++
					if done >= 100 {
						break
					}
				} else {
					other++
				}
				if other > 100 {
					break
				}
			}

			Expect(done).To(BeNumerically("==", 100))

		})

	})

})
