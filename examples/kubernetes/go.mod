module github.com/dgruber/drmaa2os/examples/kubernetes

go 1.14

replace (
	github.com/dgruber/drmaa2os => ../../../drmaa2os
	github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker => ../../../drmaa2os/pkg/jobtracker/kubernetestracker
)

require (
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.0
	github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker v0.0.0-00010101000000-000000000000
)
