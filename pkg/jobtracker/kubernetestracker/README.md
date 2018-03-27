# Kubernetes Tracker

Implements the tracker interface for kubernetes.

## Introduction

## Functionality

### Job Control Mapping

| DRMAA2 Job Control | Kubernetes      |
| :-----------------:|:---------------:|
| Suspend            | *Unsupported*   |
| Resume             | *Unsupported*   |
| Terminate          | Delete() - leads to Undetermined state |
| Hold               | *Unsupported*   |
| Release            | *Unsupported*   |

### State Mapping

Based on [JobStatus](https://kubernetes.io/docs/api-reference/batch/v1/definitions/#_v1_jobstatus)

|  DRMAA2 State.                | Kubernetes Job State  |
| :----------------------------:|:---------------------:|
| Done                          | status.Succeeded >= 1 |
| Failed                        | status.Failed >= 1    |
| Suspended                     | -                     |
| Running                       | status.Active >= 1    |
| Queued                        | -                     |
| Undetermined                  | other  / Terminate()  |


### Job Template Mapping

| DRMAA2 JobTemplate   | Kubernetes Batch Job            |
| :-------------------:|:-------------------------------:|
| RemoteCommand        | v1.Container.Command[0]         |
| Args                 | v1.Container.Args               |
| CandidateMachines[0] | v1.Container.Hostname           |
| JobCategory          | v1.Container.Image              |
| WorkingDir           | v1.Container.WorkingDir         |

Required:
* RemoteCommand
* JobCategory as it specifies the image
