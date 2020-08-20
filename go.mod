module github.com/dgruber/drmaa2os

go 1.14

//replace (
//	github.com/docker/docker => github.com/docker/docker v1.13.1
//	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
//)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/dgruber/drmaa2interface v1.0.2
	go.etcd.io/bbolt v1.3.5
)
