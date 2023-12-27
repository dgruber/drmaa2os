package dockertracker_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"

	. "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/dgruber/drmaa2interface"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job Outputter", func() {

	Context("when a job is running", func() {

		var jt drmaa2interface.JobTemplate

		var tracker *DockerTracker

		BeforeEach(func() {
			tracker, _ = New("")
			jt = drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sh",
				Args: []string{"-c", `for i in $(seq 1 100); do echo $i; sleep 0.01; done
		>&2 echo "something on stderr"
		`},
				JobCategory:    "alpine",
				StageInFiles:   map[string]string{"README.md": "/README.md"},
				JobEnvironment: map[string]string{"test": "value"},
			}
		})

		It("should be possible to get stdout and stderr", func() {

			jobID, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())

			buffer := new(bytes.Buffer)

			output := bufio.NewReadWriter(bufio.NewReaderSize(buffer, 0),
				bufio.NewWriterSize(buffer, 0))

			err = tracker.JobOutput(jobID, output)
			Ω(err).Should(BeNil())

			lineNumber := 1
			for line, _, err := output.ReadLine(); err == nil; {
				Ω(string(line)).ShouldNot(Equal(""))
				Ω(string(line)).Should(ContainSubstring(fmt.Sprintf("%d", lineNumber)))
				lineNumber++
				if lineNumber == 101 {
					// last line should be stderr output
					Ω(string(line)).Should(Equal("something on stderr"))
					break
				}
			}

		})

		It("should be possible to get stdout", func() {

			jobID, err := tracker.AddJob(jt)
			Ω(err).Should(BeNil())

			err = tracker.Wait(jobID, time.Second*30,
				drmaa2interface.Done, drmaa2interface.Failed)
			Ω(err).Should(BeNil())

			buffer := new(bytes.Buffer)

			err = tracker.JobOutput(jobID, buffer,
				JobOutputOptionNoStdError(true),
				JobOutputOptionSync(true),
				JobOutputOptionLastNLines(3))
			Ω(err).Should(BeNil())

			output := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			_, err = stdcopy.StdCopy(output, stderr, buffer)
			Ω(err).Should(BeNil())

			lastTwoLines, err := io.ReadAll(output)
			Ω(err).Should(BeNil())

			Ω(string(lastTwoLines)).Should(Equal("99\n100\n"))

		})

	})

})
