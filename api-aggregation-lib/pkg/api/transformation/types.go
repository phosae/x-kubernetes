package transformation

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Base64 generates a base64 encoding for a specified Kubernetes object.
type Base64 struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   Base64Spec
	Status Base64Status
}

type Base64Spec struct {
	// Path of the field to select in the specified object.
	// defaults to . for the entire object
	FieldPath string
}

type Base64Status struct {
	// Output is the base64-encoded representation of
	// the specified Kubernetes object or its subfield, as defined by the fieldPath.
	Output string
}
