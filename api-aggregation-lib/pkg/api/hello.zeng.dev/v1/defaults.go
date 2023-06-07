package v1

import (
	"k8s.io/apimachinery/pkg/runtime"

	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
)

func AddDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_Foo sets defaults for Foo
func SetDefaults_Foo(obj *hellov1.Foo) {
	if obj.Labels == nil {
		obj.Labels = map[string]string{}
	}
	obj.Labels["hello.zeng.dev/metadata.name"] = obj.Name
}
