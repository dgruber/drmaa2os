module github.com/dgruber/drmaa2os

go 1.15

replace (
	github.com/docker/docker => github.com/docker/docker v1.13.1
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.4 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20200413172050-18981bf12b4b
	github.com/cloudfoundry/gosigar v1.1.0
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa v0.0.0-20200831063203-2c2e9d804139
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/gorilla/mux v1.7.4
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	go.etcd.io/bbolt v1.3.5
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.18.8
)
