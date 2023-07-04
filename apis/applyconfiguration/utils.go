// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	candicepuccinicloudv1alpha1 "github.com/tliron/candice/apis/applyconfiguration/candice.puccini.cloud/v1alpha1"
	v1alpha1 "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=candice.puccini.cloud, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Device"):
		return &candicepuccinicloudv1alpha1.DeviceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DeviceDirect"):
		return &candicepuccinicloudv1alpha1.DeviceDirectApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DeviceIndirect"):
		return &candicepuccinicloudv1alpha1.DeviceIndirectApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DeviceSpec"):
		return &candicepuccinicloudv1alpha1.DeviceSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DeviceStatus"):
		return &candicepuccinicloudv1alpha1.DeviceStatusApplyConfiguration{}

	}
	return nil
}