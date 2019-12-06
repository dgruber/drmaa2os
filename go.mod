module github.com/dgruber/drmaa2os

go 1.13

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/boltdb/bolt v1.3.1
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190808214049-35bcce23fc5f
	github.com/dgruber/drmaa v0.0.0-20180514060507-3f2c6b06409b
	github.com/dgruber/drmaa2interface v1.0.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/scalingdata/ginkgo v1.1.0 // indirect
	github.com/scalingdata/go-ole v1.2.0 // indirect
	github.com/scalingdata/gomega v0.0.0-20160219221653-f331776e3035 // indirect
	github.com/scalingdata/gosigar v0.0.0-20170913211530-a501fde54c1a
	github.com/scalingdata/win v0.0.0-20150611133021-ee4771e52124 // indirect
	github.com/scalingdata/wmi v0.0.0-20170503153122-6f1e40b5b7f3 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/client-go v0.0.0-20191016111102-bec269661e48
)
