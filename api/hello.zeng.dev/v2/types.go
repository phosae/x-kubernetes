package v2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec"`
	Status FooStatus `json:"status,omitempty"`
}

type FooSpec struct {
	// Container image that the container is running to do our foo work
	Image string `json:"image"`
	// Config is the configuration used by foo container
	Config FooConfig `json:"config"`
}

type FooConfig struct {
	// Msg says hello world!
	Msg string `json:"msg"`
	// Msg1 provides some verbose information
	// +optional
	Msg1 string `json:"msg1,omitempty"`
}

// FooPhase is a label for the condition of a foo at the current time.
type FooPhase string

const (
	// FooPhaseProcessing means the pod has been accepted by the controllers, but one or more desire has not been synchorinzed
	FooPhaseProcessing FooPhase = "Processing"
	// FooPhaseReady means all conditions of foo have been meant
	FooPhaseReady FooPhase = "Ready"
)

type FooStatus struct {
	// The phase of a Foo is a simple, high-level summary of where the Foo is in its lifecycle
	// +optional
	Phase FooPhase `json:"phase,omitempty"`

	// Represents the latest available observations of a foo's current state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []FooCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type FooConditionType string

const (
	FooConditionTypeWorker FooConditionType = "Worker"
	FooConditionTypeConfig FooConditionType = "Config"
)

type FooCondition struct {
	Type   FooConditionType       `json:"type"`
	Status corev1.ConditionStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}
