### Runs the DRMAA job tracker tests in a Docker container.
test/libdrmaa:
	docker build -t drmaa/drmaajobtrackertest:latest -f ./Dockerfiles/libdrmaa/Dockerfile .
	docker run --rm -it drmaa/drmaajobtrackertest:latest

### Starts a container with Grid Engine and the latest sources for testing purposes.
libdrmaashell:
	docker build -t drmaa/drmaajobtrackertest:latest -f ./Dockerfiles/libdrmaa/Dockerfile .
	docker run --rm -it drmaa/drmaajobtrackertest:latest /bin/bash

### Runs tests the simpletracker, the job tracker for OS processes.
test/process:
	ginkgo -v pkg/jobtracker/simpletracker

### Runs docker job tracker tests.
test/docker:
	ginkgo -v pkg/jobtracker/dockertracker

### Runs Kubernetes job tracker tests.
test/kubernetes:
	ginkgo -v pkg/jobtracker/kubernetestracker

### Runs the main job tracker tests.
test: test/process test/docker test/kubernetes 

.PHONY: test/libdrmaa libdrmaashell test/process test/docker test/kubernetes test/tracker

