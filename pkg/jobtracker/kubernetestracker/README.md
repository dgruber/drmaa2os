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

_Job.Reap()_ removes the Kubernetes job and related objects from Kubernetes.

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

Using _ExtensionList_  key "env-from-secrets" will map the ":" separated secrets as
enviornment variables in the job container. The secrets must exist.

Using _ExtensionList_  key "env-from-configmaps" will map the ":" separated configmaps as
enviornment variables in the job container. The configmaps must exist.

### File staging using the Job Template

Data movement is not on core focus of the DRMAA2 standard, but it nevertheless defines two string based maps for file staging. In HPC systems data movement is usually done through parallel or network
filesystems. Cloud based systems are often using services like S3, GCS etc.

In order to simplify data movement between two pods the _StageInFiles_ and
_StageOutFiles_ maps defined in the Job Template are enhanced for smoother
Kubernetes integration. Both maps can specifiy to either move data
from the DRMAA2 process (or workflow host) to the Kubernetes jobs,
move data between two Kubernetes jobs, or transfer data back from
a Kubernetes job to the local host. Note, that some of the machanisms
have limitations by itself (like relying on Kubernetes etcd when
using ConfigMaps which has storage limits itself).

_StageInFiles_ and _StageOutFiles_  have following scheme:
- Map key specifies the target, as the target is unique.
- Map value specifies the data source.

If _StageOutFiles_ is set a sidecar (source in this project) is attached to the job
which takes care about storing the data in a persistent data structure (configmap).
This data is then downloaded the host of the DRMAA2 application.


Following source definition of _StageInFiles_ are currently implemented:
- "configmap-data:base64encodedstring" can be used to pass a byte array from the workflow process to the job. Internally a ConfigMap with the data is created in the target cluster. The ConfigMap is deleted
with the job.Reap() (Delete()) call. The ConfigMap is mounted to the file path specified in the
target definition.
- "secret-data:base64encodedstring" can be used to pass a byte array from the workflow 
process to the job. Internally a Kubernetes secret with the data is created in the 
target cluster. The Secret is removed with the job.Reap() (Delete()) call. Note, that
the content of the Secret is not stored in the JobTemplate ConfigMap in the cluster.
- "hostpath:/path/to/host/directory "can be used to mount a directory from the host where
the Kubernetes job is executed inside of the job's pod. This requires that the job
has root privileges which can be requested with the JobTemplate's extension "privileged".
- "configmap:name" Mounts an existing ConfigMap into the directory specified as target
- "pvc:name" Mounts an existing PVC with into the directory specified as target in the map. 

Target definitions of _StageInFiles_:
- /path/to/file - the path of the file in which the data from the source
definition is available (for configmap-data and secret-data)
- /path/to/directory - a path to a directory is required when using a "hostpath" 
directory as data source or a pre-existing ConfigMap or Secret.

Example:

The value in the map is the type of the volume (like _secret_ or _configmap_)
followed by a colon and the Go base64 encoded content. The key of
the map must contain the path where the file is mounted inside the job's pod.
Note that only one file can be generated with the same content as the key
in the map is unique.

Example:

    jobtemplate.StageInFiles = map[string]string{
        "/path/file.txt": "configmap-data:"+base64.StdEncoding.EncodeToString([]byte("content")),
        "/path/password.txt": "secret-data:"+base64.StdEncoding.EncodeToString([]byte("secret")),
        "/container/local/dir": "hostpath:/some/directory",
    }

### Job Template extensions

| Extension key | Extension value                   |
|:--------------|----------------------------------:|
| "labels"      | "key=value,key2=value2" v1.Labels |
| "scheduler"   | poseidon, kube-batch or any other k8s scheduler |
| "privileged"  | "true" or "TRUE"; runs container in privileged mode |

Example:

    jobtemplate.ExtensionList = map[string]string{
        "labels": "key=value",
        "privileged": "true",
    }

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

