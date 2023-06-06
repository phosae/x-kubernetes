package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func AddDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_Foo sets defaults for Foo
func SetDefaults_Foo(obj *Foo) {
	if len(obj.Spec.Msg1) == 0 {
		obj.Spec.Msg1 = obj.Spec.Msg + "🪬"
	}
	if obj.Labels == nil {
		obj.Labels = map[string]string{}
	}
	obj.Labels["hello.zeng.dev/metadata.name"] = obj.Name
}
