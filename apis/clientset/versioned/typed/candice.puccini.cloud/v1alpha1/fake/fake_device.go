// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	candicepuccinicloudv1alpha1 "github.com/tliron/candice/apis/applyconfiguration/candice.puccini.cloud/v1alpha1"
	v1alpha1 "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDevices implements DeviceInterface
type FakeDevices struct {
	Fake *FakeCandiceV1alpha1
	ns   string
}

var devicesResource = v1alpha1.SchemeGroupVersion.WithResource("devices")

var devicesKind = v1alpha1.SchemeGroupVersion.WithKind("Device")

// Get takes name of the device, and returns the corresponding device object, and an error if there is any.
func (c *FakeDevices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Device, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(devicesResource, c.ns, name), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// List takes label and field selectors, and returns the list of Devices that match those selectors.
func (c *FakeDevices) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DeviceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(devicesResource, devicesKind, c.ns, opts), &v1alpha1.DeviceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DeviceList{ListMeta: obj.(*v1alpha1.DeviceList).ListMeta}
	for _, item := range obj.(*v1alpha1.DeviceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested devices.
func (c *FakeDevices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(devicesResource, c.ns, opts))

}

// Create takes the representation of a device and creates it.  Returns the server's representation of the device, and an error, if there is any.
func (c *FakeDevices) Create(ctx context.Context, device *v1alpha1.Device, opts v1.CreateOptions) (result *v1alpha1.Device, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(devicesResource, c.ns, device), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// Update takes the representation of a device and updates it. Returns the server's representation of the device, and an error, if there is any.
func (c *FakeDevices) Update(ctx context.Context, device *v1alpha1.Device, opts v1.UpdateOptions) (result *v1alpha1.Device, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(devicesResource, c.ns, device), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDevices) UpdateStatus(ctx context.Context, device *v1alpha1.Device, opts v1.UpdateOptions) (*v1alpha1.Device, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(devicesResource, "status", c.ns, device), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// Delete takes name of the device and deletes it. Returns an error if one occurs.
func (c *FakeDevices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(devicesResource, c.ns, name, opts), &v1alpha1.Device{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDevices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(devicesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DeviceList{})
	return err
}

// Patch applies the patch and returns the patched device.
func (c *FakeDevices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Device, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(devicesResource, c.ns, name, pt, data, subresources...), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied device.
func (c *FakeDevices) Apply(ctx context.Context, device *candicepuccinicloudv1alpha1.DeviceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Device, err error) {
	if device == nil {
		return nil, fmt.Errorf("device provided to Apply must not be nil")
	}
	data, err := json.Marshal(device)
	if err != nil {
		return nil, err
	}
	name := device.Name
	if name == nil {
		return nil, fmt.Errorf("device.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(devicesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeDevices) ApplyStatus(ctx context.Context, device *candicepuccinicloudv1alpha1.DeviceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Device, err error) {
	if device == nil {
		return nil, fmt.Errorf("device provided to Apply must not be nil")
	}
	data, err := json.Marshal(device)
	if err != nil {
		return nil, err
	}
	name := device.Name
	if name == nil {
		return nil, fmt.Errorf("device.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(devicesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.Device{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Device), err
}
