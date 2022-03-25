# OS Process Tacker

## Introduction

OS Process Tracker implements the JobTracker interface used by the Go DRMAA2 implementation
in order to use standard OS processes as a backend for managing jobs as processes from the
DRMAA2 interface.

## Basic Usage

A JobTemplate requires at least:

    * RemoteCommand -> Path to the executable 

Job arrays are supported, also the control of the amount of jobs running concurrently.

### Job Control Mapping

| DRMAA2 Job Control | OS Process      |
| :-----------------:|:---------------:|
| Suspend            |  SIGTSTP        |
| Resume             |  SIGCONT        |
| Terminate          |  SIGKILL        |
| Hold               | *Unsupported*   |
| Release            | *Unsupported*   |

### State Mapping

| DRMAA2 State   | Process State       |
|:--------------:|:-------------------:|
| Queued         | *Unsupported*       |
| Running        | PID is found        |
| Suspended      |                     |
| Done           |                     |
| Failed         |                     |

### DeleteJob

Removes a finished or failed job from the internal DB to free up memory.

### Job Template Mapping

A JobTemplate is mapped into the process creation process in the following way:

| DRMAA2 JobTemplate   | OS Process                  |
| :-------------------:|:---------------------------:|
| RemoteCommand        | Executable to start         |
| JobName              |                             |
| Args                 | Arguments of the executable |
| WorkingDir           | Working directory           |
| JobEnvironment       | Environment variables set   |
| InputPath            | if set it the file that as stdin |
| OutputPath           | if set it the file that as stdout |
| ErrorPath            | if set it the file that as stderr |

JOB_ID env variable is set and TASK_ID env variable is set in case of a a job array.

### JobInfo

For finished jobs following fields could be available:

| JobInfo              | OS Process                  |
| :-------------------:|:---------------------------:|
| ExitStatus           | exit status                 |
| TerminatingSignal    | signal name                 |
| State                | Done or Failed              |
| WallclockTime        | Duration since start        |
| ID                   | process ID                  |
| AllocatedMachines    | local hostname              |
| FinishTime           | time termination is recognized |
| SubmissionHost       | local hostname              |
| JobOwner             | user ID (getuid())          |
| ExtensionList[extension.JobInfoDefaultJSessionMaxRSS]     | maxRSS |
| ExtensionList[extension.JobInfoDefaultJSessionSwap]       | nswap |
| ExtensionList[extension.JobInfoDefaultJSessionInBlock]    | inblock |
| ExtensionList[extension.JobInfoDefaultJSessionOutBlock]   | oublock |
| ExtensionList[extension.JobInfoDefaultJSessionSystemTime] | system time in ms |
| ExtensionList[extension.JobInfoDefaultJSessionUserTime]   | user time in ms |

For jobs tracked through the monitoring session following fields could be available:

| JobInfo              | OS Process                  |
| :-------------------:|:---------------------------:|
| State                | Running                     |
| DispatchTime         | Start time of process       |
| SubmissionTime       | Same as dispatch time       |
| WallclockTime        | now - dispatch time         |
| AllocatedMachines    | local hostname              |
| SubmissionHost       | local hostname              |
| JobOwner             | user ID (getuid())          |
| ExtensionList[extension.JobInfoDefaultMSessionProcessName] | process name |
| ExtensionList[extension.JobInfoDefaultMSessionCommandLine] | command line command |
| ExtensionList[jobtracker.DRMAA2_MS_JOBINFO_WORKINGDIR]     | working directory |
| ExtensionList[extension.JobInfoDefaultMSessionCPUUsage]   | how many percent of CPU time is used |
| ExtensionList[extension.JobInfoDefaultMSessionCPUAffinity] | CPU affinity list (space separated) |
| ExtensionList[extension.JobInfoDefaultMSessionMemoryUsage]   | memory usage info |
| ExtensionList[extension.JobInfoDefaultMSessionMemoryUsageRSS]   | RSS usage |
| ExtensionList[extension.JobInfoDefaultMSessionMemoryUsageVMS]   | VMS usage |
