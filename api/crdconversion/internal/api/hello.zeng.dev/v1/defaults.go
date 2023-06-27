package v1

import (
	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
)

const AnnotationImage = "spec.image"

func init() {
	localSchemeBuilder.Register(RegisterDefaults)
}

// SetDefaults_Foo sets defaults for Foo
func SetDefaults_Foo(obj *hellov1.Foo) {
	if obj.Labels == nil {
		obj.Labels = map[string]string{}
	}
	if obj.Annotations == nil {
		obj.Annotations = map[string]string{}
	}
	obj.Labels["hello.zeng.dev/metadata.name"] = obj.Name
	obj.Annotations[AnnotationImage] = "busybox:1.36"
}
