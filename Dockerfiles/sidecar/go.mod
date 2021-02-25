module github.com/dgruber/drmaa2os/Dockerfiles/sidecar

go 1.15

replace github.com/dgruber/drmaa2os/pkg/sidecar => ../../pkg/sidecar
replace k8s.io/client-go => k8s.io/client-go v0.20.2

require (
	github.com/dgruber/drmaa2os/pkg/sidecar v0.0.0-00010101000000-000000000000
	k8s.io/klog v1.0.0 // indirect
)
