package foo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	durationutil "k8s.io/apimachinery/pkg/util/duration"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
)

type fooApi struct {
	sync.RWMutex
	store map[string]*hello.Foo
}

func NewMemStore() *fooApi {
	return &fooApi{
		store: map[string]*hello.Foo{
			"default/bar": {
				ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "bar", CreationTimestamp: metav1.Now()},
				Spec: hello.FooSpec{
					Image: "busybox:1.36",
					Config: hello.FooConfig{
						Msg:  "hello world 👋",
						Msg1: "made in apiserver using k8s.io/apiserver 👊",
					},
				},
			},
		},
	}
}

var _ rest.ShortNamesProvider = &fooApi{}
var _ rest.SingularNameProvider = &fooApi{}
var _ rest.Getter = &fooApi{}
var _ rest.Lister = &fooApi{}
var _ rest.CreaterUpdater = &fooApi{}
var _ rest.GracefulDeleter = &fooApi{}
var _ rest.CollectionDeleter = &fooApi{}

// var _ rest.StandardStorage = &fooApi{} // implements all interfaces of rest.StandardStorage except rest.Watcher

func (*fooApi) ShortNames() []string {
	return []string{"fo"}
}

// GetSingularName implements rest.SingularNameProvider
func (*fooApi) GetSingularName() string {
	return hello.Resource("foo").Resource
}

// Kind implements rest.KindProvider
func (*fooApi) Kind() string {
	return "Foo"
}

// NamespaceScoped implements rest.Scoper
func (*fooApi) NamespaceScoped() bool {
	return true
}

// New implements rest.Storage
func (*fooApi) New() runtime.Object {
	return &hello.Foo{}
}

// Destroy implements rest.Storage
func (*fooApi) Destroy() {}

var simpleNameGenerator = names.SimpleNameGenerator

// Create implements rest.Creater
func (f *fooApi) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	var name, namespace, key string

	if objectMeta, err := meta.Accessor(obj); err != nil {
		return nil, errors.NewInternalError(err)
	} else {
		rest.FillObjectMetaSystemFields(objectMeta)
		if len(objectMeta.GetGenerateName()) > 0 && len(objectMeta.GetName()) == 0 {
			objectMeta.SetName(simpleNameGenerator.GenerateName(objectMeta.GetGenerateName()))
		}
		name = objectMeta.GetName()
		namespace = objectMeta.GetNamespace()
	}

	f.Lock()
	defer f.Unlock()

	key = fmt.Sprintf("%s/%s", namespace, name)
	if _, ok := f.store[key]; ok {
		return nil, errors.NewAlreadyExists(hello.Resource("foos"), key)
	}

	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}
	}

	f.store[key] = obj.(*hello.Foo)
	return obj, nil
}

// Update implements rest.Updater
func (f *fooApi) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	namespace := genericapirequest.NamespaceValue(ctx)
	key := fmt.Sprintf("%s/%s", namespace, name)

	var (
		existingObj, creatingObj runtime.Object
		creating                 = false
		err                      error
	)

	f.Lock()
	defer f.Unlock()

	if existingObj = f.store[key]; existingObj.(*hello.Foo) == nil {
		creating = true
		creatingObj = f.New()
		creatingObj, err = objInfo.UpdatedObject(ctx, creatingObj)
		if err != nil {
			return nil, false, errors.NewBadRequest(err.Error())
		}
	}

	if creating {
		creatingObj, err = f.Create(ctx, creatingObj, createValidation, nil)
		if err != nil {
			return nil, false, err
		}
		return creatingObj, true, nil
	}

	updated, err := objInfo.UpdatedObject(ctx, existingObj)
	if err != nil {
		return nil, false, errors.NewInternalError(err)
	}

	if updateValidation != nil {
		if err = updateValidation(ctx, updated, existingObj); err != nil {
			return nil, false, errors.NewBadRequest(err.Error())
		}
	}

	f.store[key] = updated.(*hello.Foo)

	return updated, false, nil
}

// Delete implements rest.GracefulDeleter
func (f *fooApi) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace := genericapirequest.NamespaceValue(ctx)
	key := fmt.Sprintf("%s/%s", namespace, name)

	f.Lock()
	defer f.Unlock()

	if obj, ok := f.store[key]; !ok { // not exists
		return nil, false, errors.NewNotFound(hello.Resource("foos"), key)
	} else {
		if deleteValidation != nil {
			if err := deleteValidation(ctx, obj); err != nil {
				return nil, false, errors.NewBadRequest(err.Error())
			}
		}
		delete(f.store, key)
		return obj, true, nil
	}
}

// DeleteCollection implements rest.CollectionDeleter
func (f *fooApi) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := genericapirequest.NamespaceValue(ctx)

	var flist hello.FooList

	f.Lock()
	defer f.Unlock()

	for _, obj := range f.store {
		if obj.Namespace == namespace {
			flist.Items = append(flist.Items, *obj)
		}
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, &flist); err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}
	}

	for _, obj := range f.store {
		if obj.Namespace == namespace {
			delete(f.store, fmt.Sprintf("%s/%s", namespace, obj.Name))
		}
	}

	return &flist, nil
}

// rest.Getter interface
func (f *fooApi) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	namespace := genericapirequest.NamespaceValue(ctx)
	key := fmt.Sprintf("%s/%s", namespace, name)

	f.RLock()
	defer f.RUnlock()

	if obj, ok := f.store[key]; !ok {
		return nil, errors.NewNotFound(hello.Resource("foos"), key)
	} else {
		return obj, nil
	}
}

// NewList implements rest.Lister
func (*fooApi) NewList() runtime.Object {
	return &hello.FooList{}
}

// List implements rest.Lister
func (f *fooApi) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	namespace := genericapirequest.NamespaceValue(ctx)

	f.RLock()
	defer f.RUnlock()

	var flist hello.FooList
	for _, obj := range f.store {
		if namespace == "" {
			flist.Items = append(flist.Items, *obj)
		} else {
			if obj.Namespace == namespace {
				flist.Items = append(flist.Items, *obj)
			}
		}
	}

	return &flist, nil
}

// ConvertToTable implements rest.Lister
func (*fooApi) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table

	table.ColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Status", Type: "string", Format: "status", Description: "status of where the Foo is in its lifecycle"},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
		{Name: "Message", Type: "string", Format: "message", Description: "foo message", Priority: 1},        // kubectl -o wide
		{Name: "Message1", Type: "string", Format: "message1", Description: "foo message plus", Priority: 1}, // kubectl -o wide
	}

	switch t := object.(type) {
	case *hello.Foo:
		table.ResourceVersion = t.ResourceVersion
		addFoosToTable(&table, *t)
	case *hello.FooList:
		table.ResourceVersion = t.ResourceVersion
		table.Continue = t.Continue
		addFoosToTable(&table, t.Items...)
	default:
	}

	return &table, nil
}

func addFoosToTable(table *metav1.Table, foos ...hello.Foo) {
	for _, foo := range foos {
		ts := "<unknown>"
		if timestamp := foo.CreationTimestamp; !timestamp.IsZero() {
			ts = durationutil.HumanDuration(time.Since(timestamp.Time))
		}
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells:  []interface{}{foo.Name, foo.Status.Phase, ts, foo.Spec.Config.Msg, foo.Spec.Config.Msg1},
			Object: runtime.RawExtension{Object: &foo},
		})
	}
}
