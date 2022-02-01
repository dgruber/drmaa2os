package main

import (
	"encoding/json"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
)

func main() {
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(
		"user", "testpassword")
	if err != nil {
		panic(err)
	}
	sm, err := drmaa2os.NewRemoteSessionManager(client.ClientTrackerParams{
		Server: "http://localhost:8088",
		Path:   "/jobserver/jobmanagement",
		Opts: []genclient.ClientOption{
			genclient.WithRequestEditorFn(basicAuthProvider.Intercept),
		},
	}, "clientjobsession.db")
	if err != nil {
		panic(err)
	}
	js, err := sm.CreateJobSession("testjobsession", "")
	if err != nil {
		// job session exists
		js, err = sm.OpenJobSession("testjobsession")
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Submitting sleep 10 job to remote server...\n")
	job, err := js.RunJob(drmaa2interface.JobTemplate{
		RemoteCommand: "sleep",
		Args:          []string{"10"},
	})
	if err != nil {
		panic(err)
	}
	err = job.WaitTerminated(drmaa2interface.InfiniteTime)
	if err != nil {
		panic(err)
	}
	jobInfo, err := job.GetJobInfo()
	if err != nil {
		panic(err)
	}
	jsonFormat, err := json.Marshal(jobInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("job info: %s\n", string(jsonFormat))
}
