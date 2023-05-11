package foo

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"

	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
)

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &hellov1.Foo{} },
		NewListFunc:               func() runtime.Object { return &hellov1.FooList{} },
		PredicateFunc:             MatchFoo,
		DefaultQualifiedResource:  hellov1.Resource("foos"),
		SingularQualifiedResource: hellov1.Resource("foos"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		TableConvertor: strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return store, nil
}
