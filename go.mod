module github.com/dgruber/drmaa2os

go 1.14

replace (
	github.com/docker/docker => github.com/docker/docker v1.13.1
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20200413172050-18981bf12b4b
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/scalingdata/ginkgo v1.1.0 // indirect
	github.com/scalingdata/go-ole v1.2.0 // indirect
	github.com/scalingdata/gomega v0.0.0-20160219221653-f331776e3035 // indirect
	github.com/scalingdata/gosigar v0.0.0-20170913211530-a501fde54c1a
	github.com/scalingdata/win v0.0.0-20150611133021-ee4771e52124 // indirect
	github.com/scalingdata/wmi v0.0.0-20170503153122-6f1e40b5b7f3 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	go.etcd.io/bbolt v1.3.5
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
)
