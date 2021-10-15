module github.com/dgruber/drmaa2os/pkg/jobtracker/libdrmaa

go 1.16

replace github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker => ../simpletracker
replace github.com/dgruber/drmaa2os/pkg/jobtracker => ../../jobtracker
replace github.com/dgruber/drmaa2os => ../../../../drmaa2os

require (
	github.com/dgruber/drmaa v1.0.0
	github.com/dgruber/drmaa2interface v1.0.2
	github.com/dgruber/drmaa2os v0.3.13
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
)
