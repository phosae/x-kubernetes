/*
boilerplate text in generated file header
*/

// Code generated by informer-gen. DO NOT EDIT.

package v2

import (
	"context"
	time "time"

	versioned "github.com/phosae/x-kubernetes/api/generated/clientset/versioned"
	internalinterfaces "github.com/phosae/x-kubernetes/api/generated/informers/externalversions/internalinterfaces"
	v2 "github.com/phosae/x-kubernetes/api/generated/listers/hello.zeng.dev/v2"
	hellozengdevv2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// FooInformer provides access to a shared informer and lister for
// Foos.
type FooInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2.FooLister
}

type fooInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFooInformer constructs a new informer for Foo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFooInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFooInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFooInformer constructs a new informer for Foo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFooInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HelloV2().Foos(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HelloV2().Foos(namespace).Watch(context.TODO(), options)
			},
		},
		&hellozengdevv2.Foo{},
		resyncPeriod,
		indexers,
	)
}

func (f *fooInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFooInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *fooInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&hellozengdevv2.Foo{}, f.defaultInformer)
}

func (f *fooInformer) Lister() v2.FooLister {
	return v2.NewFooLister(f.Informer().GetIndexer())
}
