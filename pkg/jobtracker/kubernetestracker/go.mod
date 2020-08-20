module github.com/dgruber/drmaa2os/pkg/jobtracker/kubernetestracker

go 1.14

replace (
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.9.6
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/dgruber/drmaa2os => ../../../../drmaa2os
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
)

require (
	github.com/Azure/go-autorest/autorest v0.11.4 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.2 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.2
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/googleapis/gnostic v0.4.1
	github.com/gophercloud/gophercloud v0.12.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/scalingdata/ginkgo v1.1.0 // indirect
	github.com/scalingdata/go-ole v1.2.0 // indirect
	github.com/scalingdata/gomega v0.0.0-20160219221653-f331776e3035 // indirect
	github.com/scalingdata/gosigar v0.0.0-20170913211530-a501fde54c1a // indirect
	github.com/scalingdata/win v0.0.0-20150611133021-ee4771e52124 // indirect
	github.com/scalingdata/wmi v0.0.0-20170503153122-6f1e40b5b7f3 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
)
