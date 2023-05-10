package hello_zeng_dev

import (
	"github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "hello.zeng.dev"

// SchemeGroupVersion is the hub group version for all Kinds
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

// Adds the list of known hub types to the given scheme.
// use v1 as hub type for now
//
// Normally tree of api types should like
// - hello.zeng.dev
//   - v1beta1/{types.go, register.go}
//   - v1/{types.go, register.go}
//   - types.go, register.go             <--- hub types and register
//
// all types of different kind should be translate to the internal hub type in memory
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1.Foo{},
		&v1.FooList{},
	)
	metav1.AddToGroupVersion(scheme, v1.SchemeGroupVersion)
	return nil
}
