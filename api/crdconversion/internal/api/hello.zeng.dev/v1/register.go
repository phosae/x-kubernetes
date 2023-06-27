package v1

import (
	hellov1 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "hello.zeng.dev"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	// refer and use the SchemeBuilder in api/hello.zeng.dev/v1
	// as we need add default funcs, conversion funcs...
	localSchemeBuilder = &hellov1.SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)
