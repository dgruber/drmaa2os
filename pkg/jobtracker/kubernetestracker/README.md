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
| Suspend            | _Unsupported_   |
| Resume             | _Unsupported_   |
| Terminate          | Delete() - leads to Undetermined state |
| Hold               | _Unsupported_   |
| Release            | _Unsupported_   |

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

Using _ExtensionList_  key "env-from-secrets" (or "env-from-secret") will map the ":" separated secrets listed in the map's values as enviornment variables in the job container. The secrets must exist.
(use _extension.JobTemplateK8sEnvFromSecret as key)

Using _ExtensionList_  key "env-from-configmaps" (or "env-from-configmap") will map the ":" separated configmaps listed in the map's values as enviornment variables in the job container. The configmaps must exist. (use _extension.JobTemplateK8sEnvFromConfigMap_ as key)

For more details see the [JobTemplateExtensions](#job-template-extensions)
section below.

The job's terminal output is available when the job is in a finished state
(failed or done) by the JobInfo extension key "output"
(extension.JobInfoK8sJSessionJobOutput).

```go
 if jobInfo.ExtensionList != nil {
  jobOutput, exists := jobInfo.ExtensionList[extension.JobInfoK8sJSessionJobOutput]
  if exists {
   fmt.Printf("Output of the job: %s\n", jobOutput)
  }
 }
```

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

- Map key specifies the target (in the container), as the target is unique.
- Map value specifies the data source.

Output can also be fetched through JobInfo when the job is in a terminated state.
Here the container logs are made accessible in the JobInfo extension "output".

Following source definition of _StageInFiles_ are currently implemented:

- "configmap-data:base64encodedstring" can be used to pass a byte array from the workflow
process to the job. Internally a ConfigMap with the data is created in the target cluster.
The ConfigMap is deleted with the job.Reap() (Delete()) call. The ConfigMap is mounted to
the file path specified in the target definition.
- "secret-data:base64encodedstring" can be used to pass a byte array from the workflow
process to the job. Internally a Kubernetes secret with the data is created in the
target cluster. The Secret is removed with the job.Reap() (Delete()) call. Note, that
the content of the Secret is not stored in the JobTemplate ConfigMap in the cluster.
- "hostpath:/path/to/host/directory "can be used to mount a directory from the host where
the Kubernetes job is executed inside of the job's pod. This requires that the job
has root privileges which can be requested with the JobTemplate's extension "privileged".
- "configmap:name" Mounts an existing ConfigMap into the directory specified as target
- "pvc:name" Mounts an existing PVC with into the directory specified as target in the map.
- There are more like "gce-disk", "gce-disk-read", "storageclass", "nfs" (for GoogleFilestore)...
which work similarly. Please check the convert.go file. They can also be used for staging out data
or as shared storage between multiple jobs.

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

```go
jobtemplate.StageInFiles = map[string]string{
    "/path/file.txt": "configmap-data:"+base64.StdEncoding.EncodeToString([]byte("content")),
    "/path/password.txt": "secret-data:"+base64.StdEncoding.EncodeToString([]byte("secret")),
    "/container/local/dir": "hostpath:/some/directory",
}
```

### Job Template Extensions

| Extension key | Extension value                   |
|:--------------|----------------------------------:|
| "labels"      | "key=value,key2=value2" v1.Labels |
| "scheduler"   | poseidon, kube-batch or any other k8s scheduler |
| "privileged"  | "true" or "TRUE"; runs container in privileged mode |
| "pullpolicy"  | overrides image pull policy; "always", "never", "ifnotpresent" (in any uppercase, lowercase format) |
| "distribution"  | Required for accelerators: "aks", "gke", or "eks" |
| "accelerator"  | GPU (or other type) request: "1*nvidia-tesla-v100". For aks it requires a number prefix but the type string can be arbitrary but not empty. Sets resource limits, node selector, tolerations. |
| "ttlsecondsafterfinished" | Removes the job object n seconds after it is finished. If not set the job object will never be deleted |
| "runasgroup" | Sets in the security context the GID to this number. |
| "runasuser" | Sets in the security context this UID to this number. |
| "fsgroup" | Sets in the security context this ID as filesystem group. |
| "imagepullsecrets" | Sets ImagePullSecrets in pod sec so that image can be pulled from a private registry. Comma separated list. Check out Kubernetes documentation for creating such a secret. |
| "service-account-name" | Sets the service account name for the job. |
| "node-selectors" | Sets the node selectors for the job. The format is "kubernetes.io/hostname=node1,mylabel=myvalue". |

Example:

```go
jobtemplate.ExtensionList = map[string]string{
    "labels": "key=value",
    "privileged": "true",
    "ttlsecondsafterfinished": "600",
    "imagepullsecrets": "myregistrykey",
}
```

Required for JobTemplate:

- RemoteCommand
- JobCategory as it specifies the image

Other implicit settings:

- Parallelism: 1
- Completions: 1
- BackoffLimit: 1

### Job Info Mapping

| DRMAA2 JobInfo.      | Kubernetes                           |
| :-------------------:|:------------------------------------:|
| ExitStatus           |  0 or 1 (1 if between 1 and 255 / not supported in Status)  |
| SubmissionTime       | job.CreationTimestamp.Time           |
| DispatchTime         | job.Status.StartTime.Time            |
| FinishTime           | job.Status.CompletionTime.Time       |
| State                | see above                            |
| JobID                | v1.Job.UID |
