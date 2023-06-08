/*
put install here, rather than in package hello
because hellov1 and hellov2 conversion import hello
*/
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
	hellov1 "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev/v1"
	hellov2 "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev/v2"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(hello.AddToScheme(scheme))
	utilruntime.Must(hellov1.AddToScheme(scheme))
	utilruntime.Must(hellov2.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(hellov2.SchemeGroupVersion, hellov1.SchemeGroupVersion))
}
