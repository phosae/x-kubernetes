package resource

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	durationutil "k8s.io/apimachinery/pkg/util/duration"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"

	v1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
)

var _ resource.Object = &Foo{}
var _ resource.ObjectList = &FooList{}

type Foo struct {
	v1.Foo
}

type FooList struct {
	v1.FooList
}

func (Foo) GetSingularQualifiedResource() schema.GroupResource {
	return v1.SchemeGroupVersion.WithResource("foo").GroupResource()
}

// GetGroupVersionResource implements resource.Object
func (Foo) GetGroupVersionResource() schema.GroupVersionResource {
	return v1.SchemeGroupVersion.WithResource("foos")
}

// GetObjectMeta implements resource.Object
func (f *Foo) GetObjectMeta() *metav1.ObjectMeta {
	return &f.ObjectMeta
}

// IsStorageVersion returns true -- v1.Foo is used as the internal version.
// IsStorageVersion implements resource.Object.
func (Foo) IsStorageVersion() bool {
	return true
}

// NamespaceScoped returns true to indicate Foo is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Foo) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (Foo) New() runtime.Object {
	return &v1.Foo{}
}

// NewList implements resource.Object
func (Foo) NewList() runtime.Object {
	return &v1.FooList{}
}

// GetListMeta implements resource.Object
func (c *FooList) GetListMeta() *metav1.ListMeta {
	return &c.ListMeta
}

func addFoosToTable(table *metav1.Table, foos ...v1.Foo) {
	for _, foo := range foos {
		ts := "<unknown>"
		if timestamp := foo.CreationTimestamp; !timestamp.IsZero() {
			ts = durationutil.HumanDuration(time.Since(timestamp.Time))
		}
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells:  []interface{}{foo.Name, ts, foo.Spec.Msg, foo.Spec.Msg1},
			Object: runtime.RawExtension{Object: &foo},
		})
	}
}

func (Foo) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	var table metav1.Table

	table.ColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: metav1.ObjectMeta{}.SwaggerDoc()["name"]},
		{Name: "Age", Type: "string", Description: metav1.ObjectMeta{}.SwaggerDoc()["creationTimestamp"]},
		{Name: "Message", Type: "string", Format: "message", Description: "foo message"},
		{Name: "Message1", Type: "string", Format: "message1", Description: "foo message plus", Priority: 1}, // kubectl -o wide
	}

	switch t := object.(type) {
	case *v1.Foo:
		table.ResourceVersion = t.ResourceVersion
		addFoosToTable(&table, *t)
	case *v1.FooList:
		table.ResourceVersion = t.ResourceVersion
		table.Continue = t.Continue
		addFoosToTable(&table, t.Items...)
	default:
	}

	return &table, nil
}
