package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth" // register auth plugins
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	"github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/apis/sample/v1"
	"github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/generated/openapi"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	err := builder.APIServer.
		WithOpenAPIDefinitions("sample", "v0.0.0", openapi.GetOpenAPIDefinitions).
		WithResource(&v1.Foo{}). // namespaced resource
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
