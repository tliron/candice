package controller

import (
	resources "github.com/tliron/candice/resources/candice.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *Controller) UpdateDeviceLastError(device *resources.Device, lastError string) (*resources.Device, error) {
	self.Log.Infof("updating last error to %q for device: %s/%s", lastError, device.Namespace, device.Name)

	for {
		device = device.DeepCopy()
		device.Status.LastError = lastError

		device_, err, retry := self.updateDeviceStatus(device)
		if retry {
			device = device_
		} else {
			return device_, err
		}
	}
}

func (self *Controller) updateDeviceStatus(device *resources.Device) (*resources.Device, error, bool) {
	if device_, err := self.Client.UpdateDeviceStatus(device); err == nil {
		return device_, nil, false
	} else if errors.IsConflict(err) {
		self.Log.Warningf("retrying status update for device: %s/%s", device.Namespace, device.Name)
		if device_, err := self.Client.GetDevice(device.Namespace, device.Name); err == nil {
			return device_, nil, true
		} else {
			return device, err, false
		}
	} else {
		return device, err, false
	}
}

func (self *Controller) processDevice(device *resources.Device) (bool, error) {
	return true, nil
}
