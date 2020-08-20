module github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker

go 1.14

replace (
    github.com/dgruber/drmaa2os => ../../../drmaa2os
)

require (
    github.com/googleapis/gnostic v0.4.1
	github.com/Azure/go-autorest/autorest v0.11.4 // indirect
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/utils v0.0.0-20200815180417-3bc9d57fc792 // indirect
)
