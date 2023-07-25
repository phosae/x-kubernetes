package v1beta1

import (
	transformationv1beta1 "github.com/phosae/x-kubernetes/api/transformation/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "transformation.zeng.dev"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	localSchemeBuilder = &transformationv1beta1.SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)
