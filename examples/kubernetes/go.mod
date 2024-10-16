module github.com/dgruber/drmaa2os/examples/kubernetes

go 1.21.0

toolchain go1.23.1

replace (
	github.com/dgruber/drmaa2os => ../../../drmaa2os
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
)

require (
	github.com/dgruber/drmaa2interface v1.2.1
	github.com/dgruber/drmaa2os v0.3.26
)
