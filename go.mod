module github.com/dgruber/drmaa2os

go 1.16

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.9
	github.com/containers/podman/v3 => github.com/containers/podman/v3 v3.4.4
	github.com/docker/docker => github.com/docker/docker v20.10.12+incompatible
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.4
	github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2-0.20211123152302-43a7dee1ec31
	github.com/opencontainers/runc v1.0.3 => github.com/opencontainers/runc v1.0.3
	k8s.io/api => k8s.io/api v0.22.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.6
	k8s.io/client-go => k8s.io/client-go v0.22.6
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220207220839-752842e14060
	github.com/cloudfoundry/gosigar v1.3.3
	github.com/containers/podman/v3 v3.4.4
	github.com/deepmap/oapi-codegen v1.9.1
	github.com/dgruber/drmaa v1.0.0
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/docker/docker v20.10.11+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/getkin/kin-openapi v0.89.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/ginkgo/v2 v2.0.0
	github.com/onsi/gomega v1.18.1
	github.com/opencontainers/image-spec v1.0.2
	github.com/opencontainers/runc v1.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/shirou/gopsutil/v3 v3.22.1
	go.etcd.io/bbolt v1.3.6
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	k8s.io/api v0.22.6
	k8s.io/apimachinery v0.22.6
	k8s.io/client-go v0.22.6
)
