# Tracker for Singularity Containers

_Please consider to use the process based backend (default tracker) as it can spawn any process, i.e. also singularity_

This is work in progress...ideas welcome!

## Introduction

Singularity tracker implements the JobTracker interface used by the Go DRMAA2 implementation
in order to use Singularity containers as a backend for managing jobs as containers from the
DRMAA2 interface on the same host.

## Functionality

Basically the Singularity tracker wraps the OS process tracker adding the required Singularity
calls for creating the singularity process.

## Basic Usage

A _JobTemplate_ requires at least:
  * RemoteCommand -> which is path to the executable which is started in the container
  * JobCategory -> which is the Singularity image (like vsoch-hello-world-master.simg)

If you want to see any output it makes sense to set OutputPath and ErrorPath to _/dev/stdout_
in the _JobTemplate_.

_JobTemplate_ extensions can be used to inject Singularity exec arguments like "--pid" (see _command.go_).

```go
	jt := drmaa2interface.JobTemplate{
	   RemoteCommand: "/bin/sleep",
		  Args:          []string{"600"},
		  JobCategory:   "shub://GodloveD/lolcow",
		  OutputPath:    "/dev/stdout",
		  ErrorPath:     "/dev/stderr",
	}
	// set Singularity specific arguments and options
	jt.ExtensionList = map[string]string{
		  "debug": "true",
		  "pid":   "true",
 }
```

In the _ExtensionList_ following arguments are evaluated as global singularity options:
  * debug
  * silent
  * quite
  * verbose

Boolean options are (which are injected after _singularity exec_):
  * writable
  * keep-privs
  * net
  * nv
  * overlay
  * pid
  * ipc
  * app
  * contain
  * containAll
  * userns
  * workdir
  * pwd
  * cleanenv
  * disable-cache
  * fakeroot
  * no-home
  * no-init
  * no-nv
  * no-privs
  * nohttps
  * nonet
  * rocm
  * writeable
  * writable-tmpf
  * vm
  * vm-err

Note that boolean options which are set to "false" or "FALSE" are not evaluated.

Options with values are:
  * bind
  * add-caps
  * drop-cap
  * security
  * hostname
  * network
  * network-args
  * apply-cgroups
  * scatch
  * home
  * vm-cpu
  * vm-ip
  * vm-ram
  * allow-setuid
  * app
  * drop-caps
  * app
  * pem-path

If some are missing open an issue.

# Examples

For an example please check out [singularity.go](https://github.com/dgruber/drmaa2os/blob/master/examples/singularity/singularity.go)
in the _examples_ directory.
