module github.com/phosae/x-kubernetes/api-aggregation-lib

go 1.20

require k8s.io/apimachinery v0.0.0-20230503174314-7ecc58659e5e

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20230503175222-2ef5057a4265
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20230503174314-7ecc58659e5e
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20230503185704-1ea353dbf64b
	k8s.io/client-go => k8s.io/client-go v0.0.0-20230503180226-bea472626f88
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20230503172937-f7315244e4ce
	k8s.io/component-base => k8s.io/component-base v0.0.0-20230503184328-d8237c55bb0d
	k8s.io/kms => k8s.io/kms v0.0.0-20230503185131-41fec3e2b985
)
