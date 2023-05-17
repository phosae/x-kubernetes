package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // register auth plugins
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	"github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/resource"
	"github.com/phosae/x-kubernetes/api/generated/openapi"
	"github.com/spf13/pflag"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

func main() {
	defer logs.FlushLogs()

	logOpts := logs.NewOptions()

	err := builder.APIServer.
		WithOpenAPIDefinitions("hello.zeng.dev-server", "v0.1.0", openapi.GetOpenAPIDefinitions).
		// custom backed storage, rather than default etcd
		// WithResourceAndHandler(&resource.Foo{}, func(s *runtime.Scheme, g generic.RESTOptionsGetter) (rest.Storage, error) {}).
		WithResourceAndStorage(&resource.Foo{}, func(scheme *runtime.Scheme, store *genericregistry.Store, opts *generic.StoreOptions) {
			// currently apiserver-runtime's latest update only supports k8s.io/apisver v0.26.0
			// set it SingularQualifiedResource manually, which was added at v0.27.0
			store.SingularQualifiedResource = (resource.Foo{}).GetSingularQualifiedResource()
			store.TableConvertor = (resource.Foo{})
			// replace "sigs.k8s.io/apiserver-runtime/pkg/builder/rest".GetAttrs
			// because v1.Foo doesn't implements apiserver-runtime's resource.Object 😓
			opts.AttrFunc = func(obj runtime.Object) (labels.Set, fields.Set, error) {
				accessor, ok := obj.(metav1.ObjectMetaAccessor)
				if !ok {
					return nil, nil, fmt.Errorf("given object of type %T does implements metav1.ObjectMetaAccessor", obj)
				}
				om := accessor.GetObjectMeta()
				return om.GetLabels(), fields.Set{
					"metadata.name":      om.GetName(),
					"metadata.namespace": om.GetNamespace(),
				}, nil
			}
		}).
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