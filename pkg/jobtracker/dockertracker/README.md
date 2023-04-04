# Docker Tracker

## Introduction

Docker Tracker implements the _JobTracker_ interface from Go drmaa2os.
It allows to use Docker as a backend for managing jobs as containers
from the DRMAA2 interface. It can also be used directly. The package
also contains an implementation of the _Monitorer_ interface so that
it can be used in a DRMAA2 monitoring session.

## Functionality

Docker Tracker is an API that enables the use of Docker as a backend for managing jobs as containers through the DRMAA2 interface. It provides an implementation of the _JobTracker_ interface from Go drmaa2os, allowing easy job control and management within Docker.

The functionality of Docker Tracker includes:

1. Starting Docker containers using the DRMAA2 _JobTemplate_ which requires a JobCategory (corresponding to a Docker image) and a RemoteCommand (the command to be executed within the Docker image).

2. Providing job control functions such as suspend, resume, and terminate for managing Docker containers.

3. Mapping DRMAA2 Job Control commands to corresponding Docker commands for seamless integration.

4. Mapping DRMAA2 State to Docker State to provide a consistent view of the container's status.

5. Allowing the removal of installed containers through the _DeleteJob_ command, which is equivalent to _docker rm_.

6. Supporting Job Template Mapping to efficiently map between the JobTemplate and the Docker container configuration request.

7. Implementing Job Array functionality by creating multiple tasks sequentially in a loop, since Docker does not support Array Jobs natively.

Please note that Docker Tracker does not pull container images automatically, and the required images must be pulled before using the tool. Additionally, some DRMAA2 functionalities, such as Hold and Release, are not supported in Docker Tracker due to limitations in Docker.

For the case a Docker image needs to be pulled programmatically the OS process backend can be used.

### Basic Usage

A JobTemplate requires:

    * JobCategory -> which maps to an installed Docker image
    * RemoteCommand -> which is the command executed within the given Docker image

### Job Control Mapping

| DRMAA2 Job Control | Docker          |
| :-----------------:|:---------------:|
| Suspend            | Signal: SIGSTOP |
| Resume             | Signal: SIGCONT |
| Terminate          | Signal: SIGKILL |
| Hold               | _Unsupported_   |
| Release            | _Unsupported_   |

### State Mapping

| DRMAA2 State                          | Docker State  |
| :------------------------------------:|:-------------:|
| Failed                                | OOMKilled     |
| Failed or Done depending on exit code | Exited        |
| Failed or Done depending on exit code | Dead          |
| Suspended                             | Paused        |
| Running                               | Running       |
| Queued                                | Restarting    |
| Undetermined                          | other         |

## DeleteJob

_DeleteJob_ equals _docker rm_ and is removing an installed container. It must be terminated / finished before.

### Job Template Mapping

Mapping between the job template and the Docker container config request:

| DRMAA2 JobTemplate   | Docker Container Config Request |
| :-------------------:|:-------------------------------:|
| RemoteCommand        | Cmd[0]                          |
| Args                 | Cmd[1:]                         |
| JobCategory          | Image                           |
| CandidateMachines[0] | Hostname                        |
| WorkingDir           | WorkingDir                      |
| JobEnvironment (k: v)| Env ("k=v")                     |
| StageInFiles         | -v localPath:containerPath      |
| ErrorPath            | Writes stderr into a local file (not a file in the container). |
| OutputPath           | Writes stdout into a local file (not a file in the container). |
| Extension: "user"    | User / must exist in container if set |
| Extension: "exposedPorts" | -p / multiple entries are splitted with "," |
| Extension: "net" | --net  / like "host" |
| Extension: "privileged" | --privileged  / "true"  when enabled, default "false"|
| Extension: "restart" | --restart  / like "unless-stopped", default "no" / use with care|
| Extension: "ipc" | --ipc "host" |
| Extension: "uts" | --uts "host" |
| Extension: "pid" | --pid "host" |
| Extension: "rm" | --rm  "true" or "TRUE"|

If more extensions needed just open an issue.

Note that the image must be available (pulled already)!

### Job Info Mapping

| DRMAA2 JobInfo          | Docker Container Information        |
|:-----------------------:|:-----------------------------------:|
| ID                      | Container ID                        |
| Slots                   | 1 (fixed value)                     |
| AllocatedMachines       | Config.Hostname                     |
| ExitStatus              | State.ExitCode                      |
| FinishTime              | State.FinishedAt                    |
| DispatchTime            | State.StartedAt                     |
| State                   | Mapped from Container State         |
| SubmissionTime          | Container Creation Time             |
| JobOwner                | Config.User                         |
| ExtensionList (workingdir) | Config.WorkingDir                 |
| ExtensionList (commandline) | Config.Cmd (joined as a string)   |
| ExtensionList (category) | Config.Image                      |

### Job Arrays

Since Array Jobs are not supported by Docker the job array functionality is implemented
by creating _n_ tasks sequentially in a loop. The array job ID contains all IDs of the
created Docker containers.
