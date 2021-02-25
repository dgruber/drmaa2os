module github.com/dgruber/drmaa2os

go 1.14

replace (
	github.com/docker/docker => github.com/docker/docker v1.13.1
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
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
	github.com/gophercloud/gophercloud v0.12.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	go.etcd.io/bbolt v1.3.5
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.18.8
	k8s.io/klog v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)
