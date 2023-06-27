package v2

import (
	hellov2 "github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "hello.zeng.dev"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v2"}

var (
	// refer and use the SchemeBuilder in api/hello.zeng.dev/v2
	// as we need add default funcs, conversion funcs...
	localSchemeBuilder = &hellov2.SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)


