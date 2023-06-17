package foo

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	durationutil "k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	restbuilder "sigs.k8s.io/apiserver-runtime/pkg/builder/rest"

	v2 "github.com/phosae/x-kubernetes/api-aggregation-runtime/pkg/resource/hello.zeng.dev/v2"
	hellov2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
)

type fooStorage struct {
	*genericregistry.Store
}

var _ rest.ShortNamesProvider = &fooStorage{}

func (*fooStorage) ShortNames() []string {
	return []string{"fo"}
}

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*fooStorage, error) {
	obj := v2.Foo{}
	gvr := obj.GetGroupVersionResource()
	strategy := &restbuilder.DefaultStrategy{
		Object:      obj.New(),
		ObjectTyper: scheme,
	}

	store := &genericregistry.Store{
		NewFunc:                   obj.New,
		NewListFunc:               obj.NewList,
		PredicateFunc:             strategy.Match,
		DefaultQualifiedResource:  gvr.GroupResource(),
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		StorageVersioner:          gvr.GroupVersion(),
		SingularQualifiedResource: (obj).GetSingularQualifiedResource(),
		TableConvertor:            (&fooStorage{}),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: func(obj runtime.Object) (labels.Set, fields.Set, error) {
		accessor, ok := obj.(metav1.ObjectMetaAccessor)
		if !ok {
			return nil, nil, fmt.Errorf("given object of type %T does implements metav1.ObjectMetaAccessor", obj)
		}
		om := accessor.GetObjectMeta()
		return om.GetLabels(), fields.Set{
			"metadata.name":      om.GetName(),
			"metadata.namespace": om.GetNamespace(),
		}, nil
	}}

	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &fooStorage{store}, nil
}

func addFoosToTable(table *metav1.Table, foos ...hellov2.Foo) {
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

func (*fooStorage) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table

	table.ColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Status", Type: "string", Format: "status", Description: "status of where the Foo is in its lifecycle"},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
		{Name: "Message", Type: "string", Format: "message", Description: "foo message", Priority: 1},        // kubectl -o wide
		{Name: "Message1", Type: "string", Format: "message1", Description: "foo message plus", Priority: 1}, // kubectl -o wide
	}

	switch t := object.(type) {
	case *hellov2.Foo:
		table.ResourceVersion = t.ResourceVersion
		addFoosToTable(&table, *t)
	case *hellov2.FooList:
		table.ResourceVersion = t.ResourceVersion
		table.Continue = t.Continue
		addFoosToTable(&table, t.Items...)
	default:
	}

	return &table, nil
}
