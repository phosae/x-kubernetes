package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion contains the API group and version information for the types in this package.
	SchemeGroupVersion = schema.GroupVersion{Group: "sample.zeng.dev", Version: "v1"}
	// AddToScheme applies all the stored functions to the scheme. A non-nil error
	// indicates that one function failed and the attempt was abandoned.
	AddToScheme = (&runtime.SchemeBuilder{}).AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
