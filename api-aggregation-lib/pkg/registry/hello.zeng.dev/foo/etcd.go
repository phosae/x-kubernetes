package foo

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
)

type fooStorage struct {
	*genericregistry.Store
}

var _ rest.ShortNamesProvider = &fooStorage{}

func (*fooStorage) ShortNames() []string {
	return []string{"fo"}
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*fooStorage, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &hello.Foo{} },
		NewListFunc:               func() runtime.Object { return &hello.FooList{} },
		PredicateFunc:             MatchFoo,
		DefaultQualifiedResource:  hello.Resource("foos"),
		SingularQualifiedResource: hello.Resource("foo"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		TableConvertor: strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &fooStorage{store}, nil
}
