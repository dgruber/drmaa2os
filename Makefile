# All packages which can be built without Grid Engine headers (libdrmaa)
# and without the podman v3 libraries, which do not compile anymore.
BUILDABLE_PACKAGES = $$(go list ./... | grep -v -e /libdrmaa -e podman)

### Builds and vets all packages in BUILDABLE_PACKAGES. Run it after
### updating dependencies (see the held back versions in go.mod).
vet:
	go vet $(BUILDABLE_PACKAGES)

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
	go tool ginkgo -v pkg/jobtracker/simpletracker

### Runs docker job tracker tests.
test/docker:
	go tool ginkgo -v pkg/jobtracker/dockertracker

### Runs the docker job tracker tests which do not need a Docker daemon.
test/docker/nodaemon:
	go tool ginkgo -v --label-filter='!docker' pkg/jobtracker/dockertracker

### Runs Kubernetes job tracker tests.
test/kubernetes:
	go tool ginkgo -v pkg/jobtracker/kubernetestracker

### Runs the main job tracker tests.
test: test/process test/docker test/kubernetes

.PHONY: vet test/libdrmaa libdrmaashell test/process test/docker test/docker/nodaemon test/kubernetes test/tracker

