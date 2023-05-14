package resource

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
