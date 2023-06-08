package hello

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Foo struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   FooSpec
	Status FooStatus
}

type FooSpec struct {
	// Container image that the container is running to do our foo work
	Image string
	// Config is the configuration used by foo container
	Config FooConfig
}

type FooConfig struct {
	// Msg says hello world!
	Msg string
	// Msg1 provides some verbose information
	// +optional
	Msg1 string
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
	Phase FooPhase

	// Represents the latest available observations of a foo's current state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []FooCondition
}

type FooConditionType string

const (
	FooConditionTypeWorker FooConditionType = "Worker"
	FooConditionTypeConfig FooConditionType = "Config"
)

type FooCondition struct {
	Type   FooConditionType
	Status corev1.ConditionStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FooList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Foo
}
