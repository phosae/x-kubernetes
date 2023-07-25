package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	transformation "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/transformation"
	transformationv1beta1 "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/transformation/v1beta1"
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(transformationv1beta1.AddToScheme(scheme))
	utilruntime.Must(transformation.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(transformationv1beta1.SchemeGroupVersion))
}
