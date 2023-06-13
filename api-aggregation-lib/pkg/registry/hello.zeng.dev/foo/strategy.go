package foo

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	hello "github.com/phosae/x-kubernetes/api-aggregation-lib/pkg/api/hello.zeng.dev"
)

// NewStrategy creates and returns a fooStrategy instance
func NewStrategy(typer runtime.ObjectTyper) fooStrategy {
	return fooStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Foo
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*hello.Foo)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Foo")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchFoo is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchFoo(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *hello.Foo) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type fooStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (fooStrategy) NamespaceScoped() bool {
	return true
}

func (fooStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (fooStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (fooStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	_ = obj.(*hello.Foo)
	return nil
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (fooStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string { return nil }

func (fooStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (fooStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (fooStrategy) Canonicalize(obj runtime.Object) {
}

func (fooStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (fooStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (fooStrategy) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
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
