/*
boilerplate text in generated file header
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v2

// FooConfigApplyConfiguration represents an declarative configuration of the FooConfig type for use
// with apply.
type FooConfigApplyConfiguration struct {
	Msg  *string `json:"msg,omitempty"`
	Msg1 *string `json:"msg1,omitempty"`
}

// FooConfigApplyConfiguration constructs an declarative configuration of the FooConfig type for use with
// apply.
func FooConfig() *FooConfigApplyConfiguration {
	return &FooConfigApplyConfiguration{}
}

// WithMsg sets the Msg field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Msg field is set to the value of the last call.
func (b *FooConfigApplyConfiguration) WithMsg(value string) *FooConfigApplyConfiguration {
	b.Msg = &value
	return b
}

// WithMsg1 sets the Msg1 field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Msg1 field is set to the value of the last call.
func (b *FooConfigApplyConfiguration) WithMsg1(value string) *FooConfigApplyConfiguration {
	b.Msg1 = &value
	return b
}