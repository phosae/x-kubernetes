package main

import (
	"context"
	"fmt"
	"time"

	"github.com/phosae/x-kubernetes/api/generated/clientset/versioned"
	hellov2informers "github.com/phosae/x-kubernetes/api/generated/informers/externalversions/hello.zeng.dev/v2"
	hellov2listers "github.com/phosae/x-kubernetes/api/generated/listers/hello.zeng.dev/v2"
	hellov2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, _ := clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+"/.kube/config")
	clientset := versioned.NewForConfigOrDie(config)

	// normally K8s Objects are KV based
	// client cache+indexers address this limit to some extent
	// for example we can group foos by image
	var fooInformer = hellov2informers.NewFilteredFooInformer(clientset, apiv1.NamespaceDefault, 0,
		cache.Indexers{"spec.image": func(obj interface{}) ([]string, error) {
			return []string{obj.(*hellov2.Foo).Spec.Image}, nil
		}}, nil)
	// foos, _ := podInformer.GetIndexer().ByIndex("spec.image", "busybox")

	if _, err := clientset.HelloV2().Foos(apiv1.NamespaceDefault).Create(ctx,
		&hellov2.Foo{
			ObjectMeta: metav1.ObjectMeta{Name: "myfoo"},
			Spec:       hellov2.FooSpec{Image: "busybox"},
		},
		metav1.CreateOptions{},
	); err != nil {
		panic(err)
	}
	defer func() {
		clientset.HelloV2().Foos(apiv1.NamespaceDefault).Delete(context.TODO(), "myfoo", metav1.DeleteOptions{})
	}()

	fooInformer.Run(ctx.Done())

	for !fooInformer.HasSynced() {
	}

	var fooLister = hellov2listers.NewFooLister(fooInformer.GetIndexer())
	result, err := fooLister.Foos(apiv1.NamespaceDefault).Get("myfoo")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get foo/%q from cache\n", result.GetObjectMeta().GetName())

	// using event handlers do business when Pod Object changes
	fooInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { /**/ },
		UpdateFunc: func(oldObj, newObj interface{}) { /**/ },
		DeleteFunc: func(obj interface{}) { /**/ },
	})
} //~ output: Get foo/"myfoo" from cache
