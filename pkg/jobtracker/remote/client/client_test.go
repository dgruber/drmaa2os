package client_test

import (
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dgruber/drmaa2interface"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	. "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"

	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

var _ = Describe("Client", func() {

	var testServer *httptest.Server
	var client *ClientJobTracker

	BeforeEach(func() {
		impl, _ := server.NewJobTrackerImpl(simpletracker.New("drmaa2ostestjobsession"))
		testServer = httptest.NewServer(genserver.Handler(impl))

		var err error
		client, err = New("clientdrmaa2ostestjobsession", ClientTrackerParams{
			Server: testServer.URL,
		})
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		testServer.Close()
	})

	Context("basic functionality", func() {

		It("should be able to manage a basic job lifecycle", func() {
			jobid, err := client.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			state, substate, err := client.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(substate).To(Equal(""))
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			jobs, err := client.ListJobs()
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 1))
			Expect(jobs).To(ContainElement(jobid))

			<-time.Tick(time.Millisecond * 1100)

			state, substate, err = client.JobState(jobid)
			Expect(err).To(BeNil())
			Expect(substate).To(Equal(""))
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))

			jobs, err = client.ListJobs()
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 1))
			Expect(jobs).To(ContainElement(jobid))

			err = client.DeleteJob(jobid)
			Expect(err).To(BeNil())

			jobs, err = client.ListJobs()
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 0))
		})

		It("should return the JobInfo", func() {
			jobid, err := client.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			<-time.Tick(time.Millisecond * 200)

			jobInfo, err := client.JobInfo(jobid)
			Expect(err).To(BeNil())
			Expect(jobInfo.ID).To(Equal(jobid))
			Expect(jobInfo.State.String()).To(Equal(drmaa2interface.Running.String()))
		})

		It("should list job categories", func() {
			categories, err := client.ListJobCategories()
			Expect(err).To(BeNil())
			Expect(categories).NotTo(BeNil())
		})

		It("should wait until the job is finished", func() {
			jobid, err := client.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			start := time.Now()

			err = client.Wait(jobid, time.Second*3, drmaa2interface.Done, drmaa2interface.Failed)
			Expect(err).To(BeNil())

			// it should have taken at least 1 second, the runtime of the job
			Expect(time.Now()).To(BeTemporally(">", start.Add(time.Second*1)))

			state, _, err := client.JobState(jobid)
			Expect(state.String()).To(Equal(drmaa2interface.Done.String()))
		})

	})

	Context("extensions", func() {

		It("should convert JobTemplate with extensions", func() {
			jt := drmaa2interface.JobTemplate{
				JobName: "name",
				Extension: drmaa2interface.Extension{
					ExtensionList: map[string]string{
						"extension1": "value1",
					},
				},
			}
			gen := ConvertJobTemplate(jt)
			Expect(gen.JobName).To(Equal("name"))
			Expect(gen.Extension).NotTo(BeNil())
			Expect(gen.Extension.AdditionalProperties).NotTo(BeNil())
			Expect(gen.Extension.AdditionalProperties["extension1"]).To(Equal("value1"))

			jt2 := ConvertJobTemplateToDRMAA2(gen)
			Expect(jt2.JobName).To(Equal("name"))
			Expect(jt2.ExtensionList).NotTo(BeNil())
			Expect(jt2.ExtensionList["extension1"]).To(Equal("value1"))
		})

	})

	Context("job control", func() {

		It("should be able to terminate a job", func() {
			jobid, err := client.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			err = client.JobControl(jobid, jobtracker.JobControlTerminate)
			Expect(err).To(BeNil())

			state, _, err := client.JobState(jobid)
			Expect(state.String()).To(Equal(drmaa2interface.Failed.String()))
		})

		It("should be able to suspend and release a job", func() {
			jobid, err := client.AddJob(drmaa2interface.JobTemplate{
				RemoteCommand: "/bin/sleep",
				Args:          []string{"1"},
			})
			Expect(err).To(BeNil())
			Expect(jobid).NotTo(Equal(""))

			err = client.JobControl(jobid, jobtracker.JobControlSuspend)
			Expect(err).To(BeNil())

			state, _, err := client.JobState(jobid)
			Expect(state.String()).To(Equal(drmaa2interface.Suspended.String()))

			err = client.JobControl(jobid, jobtracker.JobControlResume)
			Expect(err).To(BeNil())

			state, _, err = client.JobState(jobid)
			Expect(state.String()).To(Equal(drmaa2interface.Running.String()))

			err = client.Wait(jobid, time.Second*2, drmaa2interface.Done)
			Expect(err).To(BeNil())
		})

	})

	Context("job arrays", func() {

		It("should submit a job array and list all jobs", func() {
			arrayJobID, err := client.AddArrayJob(
				drmaa2interface.JobTemplate{
					RemoteCommand: "/bin/sleep",
					Args:          []string{"1"},
				}, 1, 5, 1, 0)
			Expect(err).To(BeNil())
			Expect(arrayJobID).NotTo(Equal(""))

			jobs, err := client.ListArrayJobs(arrayJobID)
			Expect(err).To(BeNil())
			Expect(len(jobs)).To(BeNumerically("==", 5))

			for _, jobid := range jobs {
				err := client.JobControl(jobid, jobtracker.JobControlTerminate)
				Expect(err).To(BeNil())
			}

		})

	})

})
