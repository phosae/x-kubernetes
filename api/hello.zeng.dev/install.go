package hello_zeng_dev

import (
	"github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Install registers the API group and adds types to a scheme
// DEPRECATED: implementations such as apiserver should install APIs with their own hub type, conversion func, default funcs, etc
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}
