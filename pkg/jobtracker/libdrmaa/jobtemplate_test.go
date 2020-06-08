package libdrmaa

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa"
	"github.com/dgruber/drmaa2interface"
)

var _ = Describe("Jobtemplate", func() {

	Context("basic tests", func() {

		It("should convert a JobTemplate back and forth", func() {
			originalTemplate := drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			}
			s, err := drmaa.MakeSession()
			Expect(err).To(BeNil())
			defer s.Exit()
			jt, err := s.AllocateJobTemplate()
			Expect(err).To(BeNil())
			err = ConvertDRMAA2JobTemplateToDRMAAJobTemplate(originalTemplate, &jt)
			Expect(err).To(BeNil())
			convertedJobTemplate, err := ConvertDRMAAJobTemplateToDRMAA2JobTemplate(&jt)
			Expect(err).To(BeNil())
			Expect(convertedJobTemplate.RemoteCommand).To(Equal(originalTemplate.RemoteCommand))
			Expect(len(convertedJobTemplate.Args)).To(BeNumerically("==", len(originalTemplate.Args)))
			Expect(convertedJobTemplate.Args[0]).To(Equal(originalTemplate.Args[0]))
		})

	})

	Context("Runtime tests", func() {

		It("should set the environment variables", func() {
			jt := drmaa2interface.JobTemplate{
				RemoteCommand:  "/bin/bash",
				Args:           []string{"-c", "exit $EXIT"},
				JobEnvironment: map[string]string{"EXIT": "77"},
			}
			d, err := NewDRMAATracker()
			Expect(err).To(BeNil())
			defer d.DestroySession()
			Expect(d).NotTo(BeNil())

			jobID, err := d.AddJob(jt)
			Expect(err).To(BeNil())

			err = d.Wait(jobID, time.Second*31, drmaa2interface.Failed, drmaa2interface.Done)
			Expect(err).To(BeNil())

			ji, err := d.JobInfo(jobID)
			Expect(err).To(BeNil())
			Expect(ji.ExitStatus).To(BeNumerically("==", 77))
		})

	})

})
