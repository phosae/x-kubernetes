package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Base64 generates a base64 encoding for a specified Kubernetes object.
// Base64 is commonly implemented as a subresource for specific Kubernetes kind.
// For instance, you can find it used in the
// `/apis/hello.zeng.dev/v2/namespaces/default/foos/myfoo/base64` endpoint.
type Base64 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   Base64Spec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status Base64Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type Base64Spec struct {
	// Path of the field to select in the specified object.
	// defaults to . for the entire object
	FieldPath string `json:"fieldPath,omitempty" protobuf:"bytes,1,opt,name=fieldPath"`
}

type Base64Status struct {
	// Output is the base64-encoded representation of
	// the specified Kubernetes object or its subfield, as defined by the fieldPath.
	Output string `json:"output" protobuf:"bytes,1,opt,name=output"`
}
