# Kubernetes Tracker

Implements the _JobTracker_ interface for Kubernetes batch jobs.
The _JobTracker_ is a building block for implementing the DRMAA2
interface for Kubernetes.

## Introduction

The Kubernetes tracker provides methods for managing sets of 
grouped batch jobs (within _JobSessions_). _JobSessions_ are
implemented by using labels attached to batch job objects 
("drmaa2jobsession") refering to the _JobSession_ name.

Namespaces other than "default" can be used when initializing
the _NewKubernetesSessionManager_ with _KubernetesTrackerParameters_
instead of just a ClientSet or nil.

## Functionality

## Notes

In the past Kubernetes batch jobs didn't play very well with sidecars.
So when using older Kubernetes versions and things like _istio_ you might run in state
issues (sidecar container is [still running](https://github.com/istio/istio/issues/6324)
after batch job finished).

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
| JobName              | Note: If set and a job with the same name exists in history submission will fail. metadata: Name |
| DeadlineTime         | AbsoluteTime converted to relative time (v1.Container.ActiveDeadlineSeconds) |
| JobEnvironment       | v1.EnvVar                       |


### Filestaging using the Job Template

Data movement is not on core focus of the DRMAA2 standard, but it nevertheless defines two string based maps for file staging. In 
HPC systems typically data movement is done through parallel or network
filesystems. Cloud based systems are often using services like S3, GCS etc.

In order to simplify data movement between two pods the _StageIn_ and
_StageOut_ maps defined in the Job Template are enhanced for smoother
Kubernetes integration. Both maps can specifiy to either move data
from the DRMAA2 process (or workflow host) to the Kubernetes jobs,
move data between two Kubernetes jobs, or transfer data back from
a Kubernetes job to the local host. Note, that some of the machanism
use have limitations by itself (like relying on Kubernetes etcd when
using ConfigMaps which has storage limits itself).

Both maps have following scheme:
- Map key is alwas the target, as the target is unique.
- Map value is always the source.

Following source definition of _StageInFiles_ are currently implemented:
- configmap:base64encodedstring can be used to pass a byte array from the workflow process to the job. Internally a configmap with the data is
created in the target cluster.
- secret:base64encodedstring can be used to pass a byte array from the workflow process to the job. Internally a Kubernetes secret with the data is created in the target cluster.

Target definitions of _StageInFiles_:
- /path/to/file - the path of the file in which the data from the source
definition is available

Example:

The value in the map is the type of the volume (like _secret_ or _configmap_)
followed by a colon and the Go base64 encoded content. The key of
the map must contain the path where the file is mounted inside the job's pod.
Note that only one file can be generated with the same content as the key
in the map is unique.

Example:

    jobtemplate.StageInFiles = map[string]string{
        "/path/file.txt": "configmap:"+base64.StdEncoding.EncodeToString([]byte("content")),
        "/path/password.txt": "secret:"+base64.StdEncoding.EncodeToString([]byte("secret")),
    }

Following source definition of _StageOutFiles_ are currently implemented:
- /path/to/file - the path of the file in which the data from the source
definition is read

Target definitions of _StageInFiles_:
- configmap-name:name can be used to store the data in a newly created
configmap with the name name
- /tmp/local.txt : A local path to a file which is created.

### Job Template extensions

| Extension key | Extension value                   |
|:--------------|----------------------------------:|
| "namespace"   | v1.Namespace                      |
| "labels"      | "key=value,key2=value2" v1.Labels |
| "scheduler"   | poseidon, kube-batch or any other k8s scheduler |

Required for JobTemplate:
* RemoteCommand
* JobCategory as it specifies the image

Other implicit settings:
* Parallelism: 1
* Completions: 1
* BackoffLimit: 1

### Job Info Mapping

| DRMAA2 JobInfo.      | Kubernetes                           |
| :-------------------:|:------------------------------------:|
| ExitStatus           |  0 or 1 (1 if between 1 and 255 / not supported in Status)  |
| SubmissionTime       | job.CreationTimestamp.Time           |
| DispatchTime         | job.Status.StartTime.Time            |
| FinishTime           | job.Status.CompletionTime.Time       |
| State                | see above                            |
| JobID                | v1.Job.UID |

