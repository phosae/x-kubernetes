/*
boilerplate text in generated file header
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v2

// FooSpecApplyConfiguration represents an declarative configuration of the FooSpec type for use
// with apply.
type FooSpecApplyConfiguration struct {
	Image  *string                      `json:"image,omitempty"`
	Config *FooConfigApplyConfiguration `json:"config,omitempty"`
}

// FooSpecApplyConfiguration constructs an declarative configuration of the FooSpec type for use with
// apply.
func FooSpec() *FooSpecApplyConfiguration {
	return &FooSpecApplyConfiguration{}
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *FooSpecApplyConfiguration) WithImage(value string) *FooSpecApplyConfiguration {
	b.Image = &value
	return b
}

// WithConfig sets the Config field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Config field is set to the value of the last call.
func (b *FooSpecApplyConfiguration) WithConfig(value *FooConfigApplyConfiguration) *FooSpecApplyConfiguration {
	b.Config = value
	return b
}
