package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func TestDiscovery(t *testing.T) {
	server := httptest.NewServer(BuildMux())
	defer server.Close()
	kube := kubernetes.NewForConfigOrDie(&rest.Config{Host: server.URL})

	gvrs, err := restmapper.GetAPIGroupResources(kube.DiscoveryClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*gvrs[0])

	mapper := restmapper.NewDiscoveryRESTMapper(gvrs)
	expander := restmapper.NewShortcutExpander(mapper, kube.DiscoveryClient)

	tests := []struct {
		name      string
		srcGVR    schema.GroupVersionResource
		targetGVR schema.GroupVersionResource
		targetGVK schema.GroupVersionKind
		mapperFn  func() meta.RESTMapper
		wantErr   bool
	}{
		{
			name:      `when GVR is ("","","fo") and RESTMapper is ShortcutExpander`,
			srcGVR:    schema.GroupVersionResource{Resource: "fo"},
			targetGVR: schema.GroupVersionResource{Group: "hello.zeng.dev", Version: "v1", Resource: "foos"},
			targetGVK: schema.GroupVersionKind{Group: "hello.zeng.dev", Version: "v1", Kind: "Foo"},
			mapperFn:  func() meta.RESTMapper { return expander },
			wantErr:   false,
		},
		{
			name:      `when GVR is ("","","fo") and RESTMapper is DiscoveryRESTMapper`,
			srcGVR:    schema.GroupVersionResource{Resource: "fo"},
			targetGVR: schema.GroupVersionResource{Group: "hello.zeng.dev", Version: "v1", Resource: "foos"},
			targetGVK: schema.GroupVersionKind{Group: "hello.zeng.dev", Version: "v1", Kind: "Foo"},
			mapperFn:  func() meta.RESTMapper { return mapper },
			wantErr:   true,
		},
		{
			name:      `when GVR is ("","","foo") and RESTMapper is DiscoveryRESTMapper`,
			srcGVR:    schema.GroupVersionResource{Resource: "foo"},
			targetGVR: schema.GroupVersionResource{Group: "hello.zeng.dev", Version: "v1", Resource: "foos"},
			targetGVK: schema.GroupVersionKind{Group: "hello.zeng.dev", Version: "v1", Kind: "Foo"},
			mapperFn:  func() meta.RESTMapper { return mapper },
			wantErr:   false,
		},
		{
			name:      `when GVR is ("","","foos") and RESTMapper is DiscoveryRESTMapper`,
			srcGVR:    schema.GroupVersionResource{Resource: "foos"},
			targetGVR: schema.GroupVersionResource{Group: "hello.zeng.dev", Version: "v1", Resource: "foos"},
			targetGVK: schema.GroupVersionKind{Group: "hello.zeng.dev", Version: "v1", Kind: "Foo"},
			mapperFn:  func() meta.RESTMapper { return mapper },
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gvk, err := tt.mapperFn().KindFor(tt.srcGVR)
			gvr, err1 := tt.mapperFn().ResourceFor(tt.srcGVR)
			if (err != nil) != tt.wantErr {
				t.Errorf("KindFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err1 != nil) != tt.wantErr {
				t.Errorf("ResourceFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if gvk != tt.targetGVK {
				t.Errorf("KindFor() gotGVK = %v, want %v", gvk, tt.targetGVK)
			}
			if gvr != tt.targetGVR {
				t.Errorf("KindFor() gotGVR = %v, want %v", gvk, tt.targetGVR)
			}
		})
	}
}

func TestFooAPI(t *testing.T) {
	server := httptest.NewServer(BuildMux())
	defer server.Close()
	kube := kubernetes.NewForConfigOrDie(&rest.Config{Host: server.URL})

	rest := kube.RESTClient()
	ret := rest.Verb("GET").Prefix("apis", "hello.zeng.dev", "v1").
		Namespace("default").Resource("foos").Name("bar").Do(context.Background())
	if err := ret.Error(); err != nil {
		t.Fatalf("RESTGetFoo() error = %v", err)
	}

	if err := ret.Into(&Foo{}); err != nil {
		t.Fatalf("RESTGetFoo() error = %v", err)
	}

	var newFoo = Foo{}
	b, _ := json.Marshal(&Foo{
		ObjectMeta: metav1.ObjectMeta{
			Name: "newfoo",
		},
		Spec: struct {
			Msg  string "json:\"msg\""
			Msg1 string "json:\"msg1\""
		}{
			Msg:  "hello test",
			Msg1: "new foo created from test says hi 👋",
		},
	})
	err := rest.Verb("POST").Prefix("apis", "hello.zeng.dev", "v1").
		Namespace("default").Resource("foos").Body(b).Do(context.Background()).Into(&newFoo)
	if err != nil {
		t.Fatalf("RESTCreateFoo() error = %v", err)
	}
	if newFoo.Name != "newfoo" || newFoo.Spec.Msg != "hello test" || newFoo.Spec.Msg1 != "new foo created from test says hi 👋" {
		t.Fatalf("RESTCreateFoo() error = %v, want name=%s msg=%s msg1=%s, actual %v", err, "newfoo", "hello test", "new foo created from test says hi 👋", newFoo)
	}

	err = rest.Verb("DELETE").Prefix("apis", "hello.zeng.dev", "v1").
		Namespace("default").Resource("foos").Name("newfoo").Do(context.Background()).Into(&Foo{})
	if err != nil {
		t.Fatalf("RESTDeleteFoo error = %v", err)
	}
}
