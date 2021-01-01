package client

import (
	resources "github.com/tliron/candice/resources/candice.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetDevice(namespace string, deviceName string) (*resources.Device, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if device, err := self.Candice.CandiceV1alpha1().Devices(namespace).Get(self.Context, deviceName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if device.Kind == "" {
			device = device.DeepCopy()
			device.APIVersion, device.Kind = resources.DeviceGVK.ToAPIVersionAndKind()
		}
		return device, nil
	} else {
		return nil, err
	}
}

func (self *Client) ListDevices() (*resources.DeviceList, error) {
	// TODO: all devices in cluster mode
	return self.Candice.CandiceV1alpha1().Devices(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) CreateDeviceDirect(namespace string, deviceName string, host string) (*resources.Device, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	device := &resources.Device{
		ObjectMeta: meta.ObjectMeta{
			Name:      deviceName,
			Namespace: namespace,
		},
		Spec: resources.DeviceSpec{
			Protocol: resources.DeviceProtocolNETCONF,
			Direct: &resources.DeviceDirect{
				Host: host,
			},
		},
	}

	return self.createDevice(namespace, deviceName, device)
}

func (self *Client) CreateDeviceIndirect(namespace string, deviceName string, serviceNamespace string, serviceName string, port uint64) (*resources.Device, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	device := &resources.Device{
		ObjectMeta: meta.ObjectMeta{
			Name:      deviceName,
			Namespace: namespace,
		},
		Spec: resources.DeviceSpec{
			Protocol: resources.DeviceProtocolNETCONF,
			Indirect: &resources.DeviceIndirect{
				Namespace: serviceNamespace,
				Service:   serviceName,
				Port:      port,
			},
		},
	}

	return self.createDevice(namespace, deviceName, device)
}

func (self *Client) createDevice(namespace string, deviceName string, device *resources.Device) (*resources.Device, error) {
	if device, err := self.Candice.CandiceV1alpha1().Devices(namespace).Create(self.Context, device, meta.CreateOptions{}); err == nil {
		return device, nil
	} else if errors.IsAlreadyExists(err) {
		return self.Candice.CandiceV1alpha1().Devices(namespace).Get(self.Context, deviceName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) UpdateDeviceStatus(device *resources.Device) (*resources.Device, error) {
	if device_, err := self.Candice.CandiceV1alpha1().Devices(device.Namespace).UpdateStatus(self.Context, device, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if device_.Kind == "" {
			device_ = device_.DeepCopy()
			device_.APIVersion, device_.Kind = resources.DeviceGVK.ToAPIVersionAndKind()
		}
		return device_, nil
	} else {
		return device, err
	}
}

func (self *Client) DeleteDevice(namespace string, deviceName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Candice.CandiceV1alpha1().Devices(namespace).Delete(self.Context, deviceName, meta.DeleteOptions{})
}
