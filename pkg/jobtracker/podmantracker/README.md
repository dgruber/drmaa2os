# Podman Tracker (experimental)

_Please consider to use the process based backend (default tracker) as it can spawn any process, i.e. also podman_

## Introduction

Podman Tracker implements the JobTracker interface used by the Go DRMAA2 implementation
in order to use Podman as a backend for managing jobs as containers. 

## Functionality

## Basic Usage

A JobTemplate requires:
  * JobCategory -> maps to a container image name
  * RemoteCommand -> which is the command executed within the given container image

### Job Control Mapping

| DRMAA2 Job Control | Podman          |
| :-----------------:|:---------------:|
| Suspend            | Pauses the container (does not work with rootless and cgroups v1 according to podman)|
| Resume             | Continues a paused container |
| Terminate          | Kills a running container |
| Hold               | *Unsupported*   |
| Release            | *Unsupported*   |

### State Mapping

| DRMAA2 State                          | Podman State  |
| :------------------------------------:|:-------------:|
| Failed                                | Container is not "Restarting", "Paused", "Running" had has exit code != 0|
| Done | Container is not "Restarting", "Paused", "Running", and has exit code 0|
| Suspended                             | Container is paused (if supported) |
| Running                               | Container state is "Running" or "Restarting" (sets substate "restarting") |
| Queued                                | - |
| Undetermined                          | Container inspect fails |

## DeleteJob

*DeleteJob* equals *podman rm* and is removing an installed container. It must be terminated / finished before.

### Job Template Mapping

Mapping between the job template and the Podman container config is only implemented
in the most minimalistic way. More to come.

| DRMAA2 JobTemplate   | Podman Container Config Request |
| :-------------------:|:-------------------------------:|
| RemoteCommand        | Cmd[0]                          |
| Args                 | Cmd[1:]                         |
| JobCategory          | Image                           |
| CandidateMachines[0] | spec.Hostname (container's hostname) |
| WorkingDirectory     | spec.WorkDir (working dir in the container, / if not set) |
| JobEnviornment       | spec.Env |
| ExtensionList["user"]| spec.User (user in container) |
| ExtensionList["exposedPorts] | Port forwarding in format [hostip:]hostPort:containerPort,... / spec.|
| ExtensionList["privileged"]| spec.Privileged |
| ExtensionList["restart"] | spec.RestartPolicy  / like "unless-stopped", default "no" / use with care|
| ExtensionList["ipc"] | spec.IpcNS namespace (default private) / like docker --ipc "host" |
| ExtensionList["uts"] | spec.UtsNS namespace (default private) / like docker --uts "host" |
| ExtensionList["pid"] | spec.PidNS namespace (default private) / like docker --pid "host" |
| ExtensionList["rm"] | --rm  "true" or "TRUE"|

### Job Info Mapping

### Job Arrays

Since Array Jobs are not supported by Podman the job array functionality is implemented
by creating _n_ tasks sequentially in a loop. The array job ID contains all IDs of the
created Podman containers.

