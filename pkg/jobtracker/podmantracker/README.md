# Podman Tracker (experimental)

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
| Terminate          | Stops a running container |
| Hold               | *Unsupported*   |
| Release            | *Unsupported*   |

### State Mapping

| DRMAA2 State                          | Podman State  |
| :------------------------------------:|:-------------:|
| Failed                                |      |
| Failed or Done depending on exit code |         |
| Failed or Done depending on exit code |           |
| Suspended                             |         |
| Running                               |        |
| Queued                                |     |
| Undetermined                          | other         |

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

### Job Info Mapping

### Job Arrays

Since Array Jobs are not supported by Podman the job array functionality is implemented
by creating _n_ tasks sequentially in a loop. The array job ID contains all IDs of the
created Podman containers.

