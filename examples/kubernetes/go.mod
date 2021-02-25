module github.com/dgruber/drmaa2os/examples/kubernetes

go 1.15

replace (
	github.com/dgruber/drmaa2os => ../../../drmaa2os
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

require (
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.0-beta2
)
