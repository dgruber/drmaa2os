package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
)

func main() {
	SetupHandler(simpletracker.New("jobsession"))
}

func SetupHandler(jobtracker jobtracker.JobTracker) {
	impl, _ := server.NewJobTrackerImpl(jobtracker)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        genserver.Handler(impl),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
