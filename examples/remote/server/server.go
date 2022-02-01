package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgruber/drmaa2os/pkg/jobtracker"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server"
	genserver "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/server/generated"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// run jobs behind server as simple OS processes
	jobStore, err := simpletracker.NewPersistentJobStore("job.db")
	if err != nil {
		panic(err)
	}
	processTracker, err := simpletracker.NewWithJobStore(
		"testsession", jobStore, true)
	if err != nil {
		panic(err)
	}
	RunServer(processTracker)
}

func RunServer(jobTracker jobtracker.JobTracker) {

	// connect the OpenAPI spec with the job tracker
	// interface implementation - could be anything
	impl, _ := server.NewJobTrackerImpl(jobTracker)

	// using chi router and logging + basic auth middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.BasicAuth("remotetracker",
		map[string]string{
			"user": "testpassword",
		}))

	// using the multiplexer for the case we want to serve
	// different implemenations at the same server
	m := http.NewServeMux()
	m.Handle("/jobserver/jobmanagement/",
		genserver.HandlerFromMuxWithBaseURL(
			impl, router, "/jobserver/jobmanagement"))

	// using standard http - this should be https with
	// certificates
	s := &http.Server{
		Addr:           ":8088",
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
