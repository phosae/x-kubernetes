package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"

	v1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
	v2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
)

var _ resource.Object = &Foo{}
var _ resource.MultiVersionObject = &Foo{}
var _ resource.ObjectList = &FooList{}

type Foo struct {
	v1.Foo
}

type FooList struct {
	v1.FooList
}

func TrickFrameworkInstall(s *runtime.Scheme) error {
	obj := &Foo{}
	storageVersionObj := obj.NewStorageVersionObject()
	if err := s.AddConversionFunc(&obj.Foo, storageVersionObj, func(from, to interface{}, _ conversion.Scope) error {
		return (&Foo{*from.(*v1.Foo)}).ConvertToStorageVersion(to.(runtime.Object))
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc(storageVersionObj, &obj.Foo, func(from, to interface{}, _ conversion.Scope) error {
		return (&Foo{*to.(*v1.Foo)}).ConvertFromStorageVersion(from.(runtime.Object))
	}); err != nil {
		return err
	}
	defaulter := func(in interface{}) {
		var fo *v1.Foo
		if wrap, ok := in.(*Foo); ok {
			fo = &wrap.Foo
		} else {
			fo = in.(*v1.Foo)
		}
		if fo.Labels == nil {
			fo.Labels = map[string]string{}
		}
		if fo.Annotations == nil {
			fo.Annotations = map[string]string{}
		}
		fo.Labels["hello.zeng.dev/metadata.name"] = fo.Name
		fo.Annotations[AnnotationImage] = "busybox:1.36"
	}
	s.AddTypeDefaultingFunc(obj, defaulter)
	s.AddTypeDefaultingFunc(&obj.Foo, defaulter)
	return nil
}

const AnnotationImage = "spec.image"

// ConvertFromStorageVersion implements resource.MultiVersionObject.
func (f *Foo) ConvertFromStorageVersion(storageObj runtime.Object) error {
	v2obj, ok := storageObj.(*v2.Foo)
	if !ok {
		return fmt.Errorf("storageObj is not a object of %s", "pkg/resource/hello.zeng.dev/v2.Foo")
	}
	f.ObjectMeta = v2obj.ObjectMeta
	if f.Annotations == nil {
		f.Annotations = map[string]string{}
	}
	f.Annotations[AnnotationImage] = v2obj.Spec.Image
	f.Spec.Msg = v2obj.Spec.Config.Msg
	f.Spec.Msg1 = v2obj.Spec.Config.Msg1
	return nil
}

// ConvertToStorageVersion implements resource.MultiVersionObject.
func (f *Foo) ConvertToStorageVersion(storageObj runtime.Object) error {
	v2obj, ok := storageObj.(*v2.Foo)
	if !ok {
		return fmt.Errorf("storageObj is not a object of %s", "pkg/resource/hello.zeng.dev/v2.Foo")
	}
	v2obj.ObjectMeta = f.ObjectMeta
	delete(v2obj.Annotations, AnnotationImage)
	v2obj.Spec.Image = f.Annotations[AnnotationImage]
	v2obj.Spec.Config.Msg = f.Spec.Msg
	v2obj.Spec.Config.Msg1 = f.Spec.Msg1
	return nil
}

// NewStorageVersionObject implements resource.MultiVersionObject.
func (*Foo) NewStorageVersionObject() runtime.Object {
	return &v2.Foo{}
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
	return false
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
