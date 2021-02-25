TAG?=latest
REGISTRY?=drmaa
IMAGE=drmaa2os-basic-sidecar

container: binary
	docker build --pull -t ${REGISTRY}/${IMAGE}:${TAG} .

push:
	docker -- push ${REGISTRY}/${IMAGE}:${TAG}

binary:
	GOARCH=amd64 GOOS=linux go build -a -o ./bin/sidecar ../../pkg/sidecar

.PHONY: container push binary
