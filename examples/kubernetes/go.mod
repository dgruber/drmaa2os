module github.com/dgruber/drmaa2os/examples/kubernetes

go 1.16

replace (
	github.com/dgruber/drmaa2os => ../../../drmaa2os
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
)

require (
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.26
)
