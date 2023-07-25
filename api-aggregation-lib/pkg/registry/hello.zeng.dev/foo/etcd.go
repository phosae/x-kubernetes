package foo

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
)

type REST struct {
	*genericregistry.Store
}

type fooStorage struct {
	Foo    *REST
	Config *ConfigREST
	Status *StatusREST
	Base64 *Base64REST
}

var _ rest.ShortNamesProvider = &REST{}

func (*REST) ShortNames() []string {
	return []string{"fo"}
}

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*fooStorage, error) {
	strategy := NewStrategy(scheme)
	statusStrategy := NewStatusStrategy(strategy)

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

	configStore := *store
	statusStore := *store
	statusStore.UpdateStrategy = statusStrategy
	statusStore.ResetFieldsStrategy = statusStrategy

	return &fooStorage{&REST{store}, &ConfigREST{Store: &configStore}, &StatusREST{&statusStore}, NewBase64REST(store, scheme)}, nil
}

// ConfigREST implements the config subresource for a Foo
type ConfigREST struct {
	Store *genericregistry.Store
}

var _ = rest.Patcher(&ConfigREST{})

// New creates a new Config resource
func (r *ConfigREST) New() runtime.Object {
	return &hello.Config{}
}

func (*ConfigREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *ConfigREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	fooObj, err := r.Store.Get(ctx, name, options)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewNotFound(hello.Resource("foos/config"), name)
		}
		return nil, err
	}

	foo := fooObj.(*hello.Foo)

	return configFromFoo(foo), nil
}

// Update alters the spec.config subset of an object.
// Normally option createValidation and option updateValidation are validating admission control funcs
//
//	see https://github.com/kubernetes/kubernetes/blob/d25c0a1bdb81b7a9b52abf10687d701c82704602/staging/src/k8s.io/apiserver/pkg/endpoints/handlers/patch.go#L270
//	see https://github.com/kubernetes/kubernetes/blob/d25c0a1bdb81b7a9b52abf10687d701c82704602/staging/src/k8s.io/apiserver/pkg/endpoints/handlers/update.go#L210-L216
func (r *ConfigREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	obj, _, err := r.Store.Update(
		ctx,
		name,
		&configUpdatedObjectInfo{name, objInfo},
		toConfigCreateValidation(createValidation),
		toConfigUpdateValidation(updateValidation),
		false,
		options,
	)
	if err != nil {
		return nil, false, err
	}
	foo := obj.(*hello.Foo)
	newConfig := configFromFoo(foo)
	if err != nil {
		return nil, false, errors.NewBadRequest(fmt.Sprintf("%v", err))
	}
	return newConfig, false, nil
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *ConfigREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.Store.GetResetFields()
}

// configFromFoo returns a config subresource for a foo.
func configFromFoo(foo *hello.Foo) *hello.Config {
	return &hello.Config{
		ObjectMeta: metav1.ObjectMeta{
			Name:              foo.Name,
			Namespace:         foo.Namespace,
			UID:               foo.UID,
			ResourceVersion:   foo.ResourceVersion,
			CreationTimestamp: foo.CreationTimestamp,
		},
		Spec: hello.ConfigSpec{
			Msg:  foo.Spec.Config.Msg,
			Msg1: foo.Spec.Config.Msg1,
		},
	}
}

func (r *ConfigREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.Store.ConvertToTable(ctx, object, tableOptions)
}

func toConfigCreateValidation(f rest.ValidateObjectFunc) rest.ValidateObjectFunc {
	return func(ctx context.Context, obj runtime.Object) error {
		config := configFromFoo(obj.(*hello.Foo))
		return f(ctx, config)
	}
}

func toConfigUpdateValidation(f rest.ValidateObjectUpdateFunc) rest.ValidateObjectUpdateFunc {
	return func(ctx context.Context, obj, old runtime.Object) error {
		config := configFromFoo(obj.(*hello.Foo))
		oldConfig := configFromFoo(old.(*hello.Foo))
		return f(ctx, config, oldConfig)
	}
}

// configUpdatedObjectInfo transforms existing foo -> existing config -> new config -> new foo
type configUpdatedObjectInfo struct {
	name       string
	reqObjInfo rest.UpdatedObjectInfo
}

func (c *configUpdatedObjectInfo) Preconditions() *metav1.Preconditions {
	return c.reqObjInfo.Preconditions()
}

func (c *configUpdatedObjectInfo) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	foo, ok := oldObj.DeepCopyObject().(*hello.Foo)
	if !ok {
		return nil, errors.NewBadRequest(fmt.Sprintf("expected existing object type to be Foo, got %T", foo))
	}

	// if zero-value, the existing object does not exist
	if len(foo.ResourceVersion) == 0 {
		return nil, errors.NewNotFound(hello.Resource("foos/config"), c.name)
	}

	oldConfig := configFromFoo(foo)

	// old config -> new config
	newConfigObj, err := c.reqObjInfo.UpdatedObject(ctx, oldConfig)
	if err != nil {
		return nil, err
	}
	if newConfigObj == nil {
		return nil, errors.NewBadRequest("nil update passed to Config")
	}

	config, ok := newConfigObj.(*hello.Config)
	if !ok {
		return nil, errors.NewBadRequest(fmt.Sprintf("expected input object type to be Config, but %T", newConfigObj))
	}

	// validate precondition if specified (resourceVersion matching is handled by storage)
	if len(config.UID) > 0 && config.UID != foo.UID {
		return nil, errors.NewConflict(
			hello.Resource("foos/config"),
			foo.Name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", config.UID, foo.UID),
		)
	}

	// move fields to object and return
	foo.Spec.Config.Msg = config.Spec.Msg
	foo.Spec.Config.Msg1 = config.Spec.Msg1
	foo.ResourceVersion = config.ResourceVersion

	return foo, nil
}

// StatusREST implements the REST endpoint for changing the status of a foo.
type StatusREST struct {
	store *genericregistry.Store
}

// New creates a new foo resource
func (r *StatusREST) New() runtime.Object {
	return &hello.Foo{}
}

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.store.ConvertToTable(ctx, object, tableOptions)
}
