# JobTracker Implementation on Top of _libdrmaa.so_ (Grid Engine, SLURM, ...)

_For testing using a container please call, libdrmaatest.sh in drmaa2os root directory_

Basic implementation of a _JobTracker_ wrapper for _libdrmaa.so_. This is the DRMAA version 1
c library which is shipped by many workload managers. The _JobTracker_ implementation can
be used by drmaa2os to provide a Go DRMAA2 interface for drmaa version 1. 

The _Jobtracker_ uses github.com/dgruber/drmaa Go wrapper for job submission. It supports
Grid Engine (Univa Grid Engine, SGE, Son of Grid Engine, SLURM, and more).

## Usage in drmaa2os

### Known limitations

The drmaa (v1) implementation (at least in Grid Engine) does not allow to create
different job sessions in a single process at the same point in time. The currently
also limits the _MonitoringSession_ so that no _JobSession_ can be created when a
_MonitoringSession_ is used and vice versa. If this functionality is required a
_Job_ implementation which works on command line wrappers needs to implemented.

### Default Usage

The default usage is creating a session manager which calls the _NewDRMAATracker()_
The DB is for the session manager only to store job session names etc.

    sm, err := drmaa2os.NewLibDRMAASessionManager("testdb.db")
    if err != nil {
        panic(err)
    }

### Job Persistency

If job persistency is required (like for having the jobs available after restart),
then, following initialization can be use:

    params := libdrmaa.LibDRMAASessionParams{
        ContactString:           "",
        UsePersistentJobStorage: true,
        DBFilePath:              "testdbjobs.db",
    }
    sm, err := drmaa2os.NewLibDRMAASessionManagerWithParams(params, "testdb.db")

This calls the underlying _NewDRMAATrackerWithParams()_. Contact string should be
empty unless you know what you are doing. If _UsePersistentJobStorage_ is turned
on the _DBFilePath_ must be specified in which job related information is written.
If the DB file does not exist it will be created. The contact string of the underlying
drmaa1 session is written in the session manager DB, and when re-connecting to the
same session name it transparently uses it. Hence still running jobs can be still
available after application restart.

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
| ExtensionList map["DRMAA1_NATIVE_SPECIFICATION"]value | SetNativeSpecification("value")|

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