module github.com/monzo/egress-operator

go 1.13

require (
	github.com/envoyproxy/go-control-plane v0.9.2-0.20191209214638-ccc55b462344
	github.com/go-logr/logr v0.1.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	golang.org/x/sys v0.0.0-20220817070843-5a390386f1f2 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.0
	sigs.k8s.io/yaml v1.1.0
)
