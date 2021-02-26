module github.com/dgruber/drmaa2os

go 1.15

replace (
	github.com/docker/docker => github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	//github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
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
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa v0.0.0-20210226091710-ceb83e9b4fff
	github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/gorilla/mux v1.8.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1
	go.etcd.io/bbolt v1.3.5
	golang.org/x/net v0.0.0-20210226101413-39120d07d75e
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)
