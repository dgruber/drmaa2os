module github.com/dgruber/drmaa2os

go 1.16

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.8
	github.com/containers/podman/v3 => github.com/containers/podman/v3 v3.4.4
	github.com/docker/docker => github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.4
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20201123235753-4f46d6348a05
	github.com/cloudfoundry/gosigar v1.1.0
	github.com/containers/podman/v3 v3.2.0
	github.com/deepmap/oapi-codegen v1.8.1
	github.com/dgruber/drmaa v1.0.0
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/docker/docker v20.10.11+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/getkin/kin-openapi v0.61.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/opencontainers/image-spec v1.0.2-0.20210819154149-5ad6f50d6283
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.6
	golang.org/x/net v0.0.0-20211005001312-d4b1ae081e3b
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.20.6
	sigs.k8s.io/structured-merge-diff/v4 v4.0.3 // indirect
)
