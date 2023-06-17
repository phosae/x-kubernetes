package main

import (
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // register auth plugins
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	fooregistry "github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/registry/hello/foo"
	hellov1 "github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/resource/hello.zeng.dev/v1"
	hellov2 "github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/resource/hello.zeng.dev/v2"
	"github.com/phosae/x-kubernetes/api/generated/openapi"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

func main() {
	defer logs.FlushLogs()

	logOpts := logs.NewOptions()

	err := builder.APIServer.
		WithOpenAPIDefinitions("hello.zeng.dev-server", "v0.1.0", openapi.GetOpenAPIDefinitions).
		// Go can't implement interface for external package api/hello.zeng.dev/{version}
		// So manually registry conversion for external types to trick apiserver runtime framework
		WithAdditionalSchemeInstallers(func(s *runtime.Scheme) error {
			return hellov1.TrickFrameworkInstall(s)
		}).
		// customize backed storage (can be replace with any implemention instead of etcd
		// normally use WithResourceAndStorage is ok
		// we choose WithResourceAndHandler only because WithResourceAndStorage don't support shortNames
		WithResourceAndHandler(&hellov1.Foo{}, func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
			return fooregistry.NewREST(scheme, optsGetter)
		}).
		WithResourceAndHandler(&hellov2.Foo{}, func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
			return fooregistry.NewREST(scheme, optsGetter)
		}).
		// WithResourceAndStorage(&resource.Foo{}, func(scheme *runtime.Scheme, store *genericregistry.Store, opts *generic.StoreOptions) {
		// 	// currently apiserver-runtime's latest update only supports k8s.io/apisver v0.26.0
		// 	// set it SingularQualifiedResource manually, which was added at v0.27.0
		// 	store.SingularQualifiedResource = (resource.Foo{}).GetSingularQualifiedResource()
		// 	store.TableConvertor = (resource.Foo{})
		// 	// replace "sigs.k8s.io/apiserver-runtime/pkg/builder/rest".GetAttrs
		// 	// because v1.Foo doesn't implements apiserver-runtime's resource.Object 😓
		// 	opts.AttrFunc = func(obj runtime.Object) (labels.Set, fields.Set, error) {
		// 		accessor, ok := obj.(metav1.ObjectMetaAccessor)
		// 		if !ok {
		// 			return nil, nil, fmt.Errorf("given object of type %T does implements metav1.ObjectMetaAccessor", obj)
		// 		}
		// 		om := accessor.GetObjectMeta()
		// 		return om.GetLabels(), fields.Set{
		// 			"metadata.name":      om.GetName(),
		// 			"metadata.namespace": om.GetNamespace(),
		// 		}, nil
		// 	}
		// }).
		WithOptionsFns(func(so *builder.ServerOptions) *builder.ServerOptions {
			// do log opts trick
			logs.InitLogs()
			logsapi.ValidateAndApply(logOpts, utilfeature.DefaultFeatureGate)
			return so
		}).
		WithFlagFns(func(ss *pflag.FlagSet) *pflag.FlagSet {
			logsapi.AddFlags(logOpts, ss)
			return ss
		}).
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
