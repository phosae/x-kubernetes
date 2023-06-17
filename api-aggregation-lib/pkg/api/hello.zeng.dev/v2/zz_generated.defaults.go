//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
boilerplate text in generated file header
*/

// Code generated by defaulter-gen. DO NOT EDIT.

package v2

import (
	v2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v2.Foo{}, func(obj interface{}) { SetObjectDefaults_Foo(obj.(*v2.Foo)) })
	scheme.AddTypeDefaultingFunc(&v2.FooList{}, func(obj interface{}) { SetObjectDefaults_FooList(obj.(*v2.FooList)) })
	return nil
}

func SetObjectDefaults_Foo(in *v2.Foo) {
	SetDefaults_Foo(in)
}

func SetObjectDefaults_FooList(in *v2.FooList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_Foo(a)
	}
}