// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// DeviceIndirectApplyConfiguration represents an declarative configuration of the DeviceIndirect type for use
// with apply.
type DeviceIndirectApplyConfiguration struct {
	Namespace *string `json:"namespace,omitempty"`
	Service   *string `json:"service,omitempty"`
	Port      *uint64 `json:"port,omitempty"`
}

// DeviceIndirectApplyConfiguration constructs an declarative configuration of the DeviceIndirect type for use with
// apply.
func DeviceIndirect() *DeviceIndirectApplyConfiguration {
	return &DeviceIndirectApplyConfiguration{}
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *DeviceIndirectApplyConfiguration) WithNamespace(value string) *DeviceIndirectApplyConfiguration {
	b.Namespace = &value
	return b
}

// WithService sets the Service field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Service field is set to the value of the last call.
func (b *DeviceIndirectApplyConfiguration) WithService(value string) *DeviceIndirectApplyConfiguration {
	b.Service = &value
	return b
}

// WithPort sets the Port field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Port field is set to the value of the last call.
func (b *DeviceIndirectApplyConfiguration) WithPort(value uint64) *DeviceIndirectApplyConfiguration {
	b.Port = &value
	return b
}