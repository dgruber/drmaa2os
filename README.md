# drmaa2os - A Go API for OS Processes, Docker Containers, Cloud Foundry Tasks, Kubernetes Jobs, Grid Engine Jobs, Podman containers, and more...

_DRMAA2 for OS processes and more_

[![CircleCI](https://circleci.com/gh/dgruber/drmaa2os.svg?style=svg)](https://circleci.com/gh/dgruber/drmaa2os)
[![codecov](https://codecov.io/gh/dgruber/drmaa2os/branch/master/graph/badge.svg)](https://codecov.io/gh/dgruber/drmaa2os)

> _Update_: The Go DRMAA2 interface and the implementation based on the JobTracker
> interface are now decoupled. In order to use a specific backend, like Docker,
> the package providing the JobTracker implementation needs to be imported so
> that the init() method is called for registering at the DRMAA2 implementation.
> 
> Like when using the Docker backend:
> 
> ```
> 	_ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
> ```


This is a Go API based on an open standard ([Open Grid Forum DRMAA2](https://www.ogf.org/documents/GFD.231.pdf)) for submitting and
supervising workloads running in operating system processes, containers, PODs, tasks, or HPC batch jobs.

The API allows you to develop and run job workflows in OS processes and switch later to 
containers running in Kubernetes, as Cloud Foundry tasks, pure Docker, Singularity, 
or any HPC workload manager which supports the DRMAA standard through the C _libdrmaa.so_
library (like SLURM, Grid Engine, ...) without changing the application logic.

Its main pupose is supporting you with an abstraction layer on top of platforms, workload managers, 
and HPC cluster schedulers, so that a software developer don't need to deal with the underlaying details and differences of job submission, status checking, and more.

An even simpler interface for creating job workflows without dealing with the DRMAA2 details is
[*wfl*](https://github.com/dgruber/wfl) which is based on the Go DRMAA2 implementation.

For details about the mapping of job operations please consult the platform specific READMEs:

  * [OS Processes](pkg/jobtracker/simpletracker/README.md)
  * [Cloud Foundry](pkg/jobtracker/cftracker/README.md)
  * [Docker / Moby](pkg/jobtracker/dockertracker/README.md)
  * [Kubernetes](pkg/jobtracker/kubernetestracker/README.md)
  * [Singularity](pkg/jobtracker/singularity/README.md)
  * [libdrmaa.so](pkg/jobtracker/libdrmaa/README.md)
  * [Podman](pkg/jobtracker/podmantracker/README.md)

[Feedback](mailto:info@gridengine.eu) welcome!

For a Go DRMAA2 wrapper based on C DRMAA2 (_libdrmaa2.so_) like for *Univa Grid Engine* please
check out [drmaa2](https://github.com/dgruber/drmaa2).

## Basic Usage

Following example demonstrates how a job running as OS process can be executed. More examples can be found in the _examples_ subdirectory.
Per default jobs are managed in main memory hence after restarting your app all processes are not visible to your app even they
are running. If persistency between restarts is required, please use _NewDefaultSessionManagerWithParams()_ with a
_simpletracker.SimpleTrackerInitParams_ as argument.

Note that at this point in time only _JobSessions_ are implemented.

```go
    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
    )

	sm, err := drmaa2os.NewDefaultSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		panic(err)
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "sleep",
		Args:          []string{"2"},
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	if job.GetState() == drmaa2interface.Done {
		job2, _ := js.RunJob(jt)
		job2.WaitTerminated(drmaa2interface.InfiniteTime)
	} else {
		fmt.Println("Failed to execute job1 successfully")
	}

	js.Close()
	sm.DestroyJobSession("jobsession")
```

## Using other Backends

Using other backends for workload management and execution only differs in creating
a different _SessionManager_. Different _JobTemplate_ attributes might be neccessary when
switching the implementation. If using a backend which supports container images it
might be required to set the _JobCategory_ to the container image name.

### Docker

If Docker is installed locally it will automatically detect it. For pointing to
a different host environment variables needs to be set before the _SessionManager_
is created.

"Use DOCKER_HOST to set the url to the docker server.
 Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
 Use DOCKER_CERT_PATH to load the TLS certificates from.
 Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default."

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/dockertracker"
    )

	sm, err := drmaa2os.NewDockerSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		panic(err)
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "sleep",
		Args:          []string{"2"},
		JobCategory:   "busybox",
	}
	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	js.Close()
	sm.DestroyJobSession("jobsession")
```

### Kubernetes

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker"
    )

	sm, err := drmaa2os.NewKubernetesSessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		panic(err)
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "sleep",
		Args:          []string{"2"},
		JobCategory:   "busybox",
	}
	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	js.Close()
	sm.DestroyJobSession("jobsession")
```

### Cloud Foundry

The Cloud Foundry _SessionManager_ requires details for connecting to the
Cloud Foundry cloud controller API when being created. The _JobCategory_ needs to
be set to the application GUID which is the source of the container image
of the task.

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/cftracker"
    )

	sm, err := drmaa2os.NewCloudFoundrySessionManager("api.run.pivotal.io", "user", "password", "test.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		panic(err)
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "dbbackup.sh",
		Args:          []string{"location"},
		JobCategory:   "123CFAPPGUID",
	}
	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	js.Close()
	sm.DestroyJobSession("jobsession")
```

### Singularity

The Singularity _SessionManager_ wraps the singularity command which is required to be installed.
The container images can be provided in any form (like pointing to file or shub) but are 
required to be set as _JobCategory_ for each job.

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/singularity"
    )

	sm, err := drmaa2os.NewSingularitySessionManager("testdb.db")
	if err != nil {
		panic(err)
	}

	js, err := sm.CreateJobSession("jobsession", "")
	if err != nil {
		panic(err)
	}

	jt := drmaa2interface.JobTemplate{
		RemoteCommand: "sleep",
		Args:          []string{"2"},
		JobCategory:   "shub://GodloveD/lolcow",
	}

	job, err := js.RunJob(jt)
	if err != nil {
		panic(err)
	}

	job.WaitTerminated(drmaa2interface.InfiniteTime)

	js.Close()
	sm.DestroyJobSession("jobsession")
```

### DRMAA (version 1) - libdrmaa.so

The _LibDRMAASessionManager_ can be used for submitting jobs through a pre-existing _libdrmaa.so_ which
is available and supported by many HPC workload managers (like Univa Grid Engine, SLURM, PBS, LSF,
Son of Grid Engine, ...).

There are a few things to consider at compile time and runtime. The CGO_LDFLAGS and CGO_CFLAGS must be
set according to the documentation in [https://github.com/dgruber/drmaa](https://github.com/dgruber/drmaa).
Also the LD_LIBRARY_PATH needs to be set accordingly.

An example using Grid Engine running in a container is [here](https://github.com/dgruber/drmaa2os/tree/master/examples/libdrmaa)

The compile time configuration is external meaning the C library must be in the path or LD_LIBRARY_PATH and 
CGO_LDFLAGS and CGO_CFLAGS must be set according to the documentation in [https://github.com/dgruber/drmaa](https://github.com/dgruber/drmaa).

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa"
    )
    
	sm, err := drmaa2os.NewLibDRMAASessionManager("testdb.db")
	if err != nil {
		panic(err)
	}
```

### Podman (Remote)

First experimental version is implemented and tested on macos accessing Podman
on a remote VM. When compiling on macos _brew install gpgme_ helped me getting
the C header dependencies of Podman installed. Accessing podman can be achieved
through _ssh_ in that case (calling podman system service --time=0 unix:///tmp/podman.sock
in the podman VM for which the ssh port is defined at localhost:2222 on a Vagrant
based vbox VM).

If _ConnectionURIOverride_ is not set the implementation uses the default connection
to the Podman REST API server. This server can be setup by _podman system service -t 0 &_
in Linux enviornments. 

Note, that it currently the implementation expects that the images are pre-pulled.

For running podman locally the process based implementation (simpletracker) can
be used.

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/podmantracker"
    )
    
	sm, err := drmaa2os.NewPodmanSessionManager(PodmanTrackerParams{
					ConnectionURIOverride: "ssh://vagrant@localhost:2222/tmp/podman.sock?secure=False",
				}, "testdb.db")
	if err != nil {
		panic(err)
	}
```

### Remote

The _remote_ directory in _/pkg/jobtracker_ contains a client/server implementation of the
_JobTracker_ interface allowing to create clients and server for any backends (_JobTracker_ 
implementations) mentioned above. The client/server protocol is defined in OpenAPI v3. Based
on that _Go_ client and server stubs have been generated using _oapi-codegen_. The OpenAPI
spec contains also the DRMAA2 data types which might be useful for other projects.

The remote _JobTracker_ server can be used in any Go DRMAA2 application.

```go

    import (
        "github.com/dgruber/drmaa2os
        _ "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"
    )
    
	sm, err := drmaa2os.NewRemoteSessionManager(ClientTrackerParams{
					Server: "localhost:8080",
				}, "testdb.db")
	if err != nil {
		panic(err)
	}
```

The server can be implemented by using any JobTracker implementation as
argument in the server implementation.

```go

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
```


