# Containerd Tracker

## Introduction

Containerd Tracker provides an implementation of the [DRMAA2 JobTracker](https://github.com/dgruber/drmaa2interface) interface for managing containers with [containerd](https://containerd.io/). It allows you to use the DRMAA2 interface to create, control, and manage containerd containers as jobs. The package also contains implementations for managing container lifecycles, getting container information, and mapping container states to DRMAA2 job states.

## Functionality

Containerd Tracker provides an API to manage containerd containers as jobs using the DRMAA2 interface. It offers the following functionality:

1. Creating and starting containerd containers using the DRMAA2 JobTemplate. The container image and command to run within the container are specified in the JobCategory and Args fields of the JobTemplate, respectively.

2. Listing all containerd containers and their corresponding job IDs.

3. Providing job control functions such as suspend, resume, and terminate for managing containerd containers.

4. Mapping DRMAA2 Job Control commands to corresponding containerd actions for seamless integration.

5. Mapping DRMAA2 JobState to containerd state to provide a consistent view of the container's status.

6. Retrieving container information and mapping it to the DRMAA2 JobInfo struct.

Please note that some DRMAA2 features, such as Hold, Release, and Job Arrays, are not supported in the Containerd Tracker due to limitations in containerd.

## Basic Usage

To use Containerd Tracker, you'll need to create a new ContainerdJobTracker instance with the containerd address:

```go
tracker, err := containerdtracker.NewContainerdJobTracker("/run/containerd/containerd.sock")
```

A JobTemplate requires:

* JobCategory -> which maps to the container image to be used
* Args -> which is the command to be executed within the given container image

### Job Control Mapping

| DRMAA2 Job Control | Containerd Action |
|:------------------:|:-----------------:|
| Suspend            | Pause             |
| Resume             | Resume            |
| Terminate          | Kill (SIGTERM)    |

### State Mapping

| DRMAA2 State        | Containerd State  |
|:-------------------:|:-----------------:|
| Queued              | Created               |
| Running             | Running, Pausing      |
| Suspended           | Paused                |
| Done                | Stopped - Exit code 0 |
| Failed              | Stopped - Exit code != 0 |
| Undetermined        | other                 |

## Job Info

The Containerd Tracker provides an implementation to retrieve container information and map it to the DRMAA2 JobInfo struct. Some fields may not be directly available from containerd and might require further customization based on your specific requirements.

### Job Arrays

Job Arrays are not supported in the Containerd Tracker due to limitations in containerd. However, you can implement job arrays by manually creating multiple tasks sequentially in a loop.

## Testing

The Containerd Tracker includes Ginkgo tests that can be run to ensure proper functionality. To run the tests, make sure you have Ginkgo and Gomega installed, and then execute `ginkgo` in the package directory.

## Limitations

The Containerd Tracker has some limitations, including:

* It does not support DRMAA2 Hold and Release actions.
* Job Arrays are not natively supported in containerd and must be implemented manually.
* Some JobTemplate fields may require additional customization to work seamlessly with containerd configurations.

Despite these limitations, the Containerd Tracker provides a convenient way to use the DRMAA2 interface for managing containerd containers as jobs.