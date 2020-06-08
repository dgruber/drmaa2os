# JobTracker implementation on top of libdrmaa.so

Basic implementation of a _JobTracker_ wrapper for _libdrmaa.so_. This is the DRMAA version 1
c library which is shipped with many workload managers. The JobTracker implementation can
be used by DRMAA2OS to provide a Go DRMAA2 interface for drmaa version 1. 

It is currently tested with Grid Engine using the Docker image in this directory.
The _Jobtracker_ uses github.com/dgruber/drmaa Go wrapper for job submission.

## JobTemplate Mapping

| DRMAA2 JobTemplate | Internal Go drmaa job template  |
|---|---|
| RemoteCommand  | SetRemoteCommand |
| Args  | SetArgs  |
| InputPath | SetInputPath(":"+InputPath) |
| OutputPath | SetOutputPath(":"+OutputPath) |
| ErrorPath | SetErrorPath(":"+ErrorPath) |
| JobName | "cdrmaatrackerjob" if not set / SetJobName |
| JoinFiles | SetJoinFiles |
| Email | SetEmail |
| JobEnviornment map[key]value | SetJobEnviornment("key=value", ...)|

## JobState Mapping

The following table shows how DRMAA2 job states are mapped to DRMAA version 1
job states.

| DRMAA2 Job State | Internal Go drmaa job state |
|---|---|
| drmaa2interface.Undetermined | drmaa.PsUndetermined |
| drmaa2interface.Queued | drmaa.PsQueuedActive |
| drmaa2interface.QueuedHeld | drmaa.PsSystemOnHold |
| drmaa2interface.QueuedHeld | drmaa.PsUserOnHold |
| drmaa2interface.QueuedHeld | drmaa.PsUserSystemOnHold |
| drmaa2interface.Running | drmaa.PsRunning |
| drmaa2interface.Suspended  | drmaa.PsSystemSuspended |
| drmaa2interface.Suspended | drmaa.PsUserSuspended |
| drmaa2interface.Suspended | drmaa.PsUserSystemSuspended |
| drmaa2interface.Done | drmaa.PsDone |
| drmaa2interface.Failed | drmaa.PsFailed |


## JobInfo Mapping

_JobInfo_ is set when the job is in an end state.

| DRMAA2 JobInfo | Internal Go drmaa job info |
|---|---|
| ExitStatus | Only meaningful when drmaa job _HasExited()_  |
| ID  | JobID()  |
| SubmissionTime | Resource usage map submission_time value |
| DispatchTime | Resource usage map start_time value |
| FinishTime | Resource usage map end_time value |

_TODO_ add more