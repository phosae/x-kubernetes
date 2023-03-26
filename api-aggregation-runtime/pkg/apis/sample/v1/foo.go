package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
)

var _ resource.Object = &Foo{}
var _ resource.ObjectList = &FooList{}

// +genclient:namespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a demo type integrated into apiserver runtime framework.
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FooSpec `json:"spec"`
}

type FooSpec struct {
	// Msg persist things like `hello world`
	Msg string `json:"msg"`
}

// +genclient:namespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo objects.
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

// GetGroupVersionResource implements resource.Object
func (Foo) GetGroupVersionResource() schema.GroupVersionResource {
	return SchemeGroupVersion.WithResource("foos")
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

// NamespaceScoped returns false to indicate Fischer is NOT a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Foo) NamespaceScoped() bool {
	return false
}

// New implements resource.Object
func (Foo) New() runtime.Object {
	return &Foo{}
}

// NewList implements resource.Object
func (Foo) NewList() runtime.Object {
	return &FooList{}
}

// GetListMeta implements resource.Object
func (c *FooList) GetListMeta() *metav1.ListMeta {
	return &c.ListMeta
}
