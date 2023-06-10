package main

import (
	"context"
	"fmt"

	applyhellov2 "github.com/phosae/x-kubernetes/api/generated/applyconfiguration/hello.zeng.dev/v2"
	"github.com/phosae/x-kubernetes/api/generated/clientset/versioned"
	hellov2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, _ := clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+"/.kube/config")
	clientset, _ := versioned.NewForConfig(config)

	fooClient := clientset.HelloV2().Foos(apiv1.NamespaceDefault)
	fo := &hellov2.Foo{
		ObjectMeta: metav1.ObjectMeta{Name: "myfoo", Labels: map[string]string{"app": "myfoo"}},
		Spec:       hellov2.FooSpec{Image: "busybox", Config: hellov2.FooConfig{Msg: "hi 👋"}},
	}

	result, err := fooClient.Create(context.TODO(), fo, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created foo %q.\n", result.GetObjectMeta().GetName())
	defer func() { fooClient.Delete(context.TODO(), fo.Name, metav1.DeleteOptions{}) }()

	// similar to: kubectl patch fo/myfoo -p '{"spec": {"config": {"msg1": "hello world 🔮"}}}'
	fooApplyCfg := applyhellov2.Foo("myfoo", apiv1.NamespaceDefault).
		WithSpec(applyhellov2.FooSpec().WithConfig(applyhellov2.FooConfig().WithMsg1("hello world 🔮")))
	applyRet, err := fooClient.Apply(context.Background(), fooApplyCfg, metav1.ApplyOptions{FieldManager: "example", Force: false})
	if err != nil {
		panic(err)
	}
	fmt.Printf("foo/%s Config:\n", applyRet.Name)
	fmt.Printf("\tmsg: %v\n", applyRet.Spec.Config.Msg)
	fmt.Printf("\tmsg1: %v\n", applyRet.Spec.Config.Msg1)
} //~ ouput:
/*
Created foo "myfoo".
foo/myfoo Config:
        msg: hi 👋
        msg1: hello world 🔮
*/
