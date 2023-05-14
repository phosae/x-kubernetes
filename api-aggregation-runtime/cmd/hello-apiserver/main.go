package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // register auth plugins
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	"github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/resource"
	"github.com/phosae/x-kubernetes/api/generated/openapi"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	err := builder.APIServer.
		WithOpenAPIDefinitions("hello.zeng.dev-server", "v0.1.0", openapi.GetOpenAPIDefinitions).
		// custom backed storage, rather than default etcd
		// WithResourceAndHandler(&resource.Foo{}, func(s *runtime.Scheme, g generic.RESTOptionsGetter) (rest.Storage, error) {}).
		WithResourceAndStorage(&resource.Foo{}, func(scheme *runtime.Scheme, store *genericregistry.Store, opts *generic.StoreOptions) {
			// currently apiserver-runtime's latest update only supports k8s.io/apisver v0.26.0
			// set it SingularQualifiedResource manually, which was added at v0.27.0
			store.SingularQualifiedResource = (resource.Foo{}).GetSingularQualifiedResource()
		}).
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
