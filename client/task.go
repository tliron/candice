package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	resources "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/transcribe"
)

var permissions int64 = 0700

func (self *Client) GetTask(namespace string, componentName string, taskName string) (string, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, componentName, taskName)
	if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		for _, podName := range podNames {
			self.Log.Infof("setting task %q in operator pod: %s/%s", taskName, self.Namespace, podName)
			var buffer bytes.Buffer
			if self.ReadFromContainer(self.Namespace, podName, "operator", &buffer, taskPath); err == nil {
				return buffer.String(), nil
			}
		}
	} else {
		return "", err
	}

	// TODO: not found error
	return "", nil
}

func (self *Client) SetTask(namespace string, componentName string, taskName string, content string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, componentName, taskName)
	if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		for _, podName := range podNames {
			self.Log.Infof("setting task %q in operator pod: %s/%s", taskName, self.Namespace, podName)
			if err := self.WriteToContainer(self.Namespace, podName, "operator", strings.NewReader(content), taskPath, &permissions); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func (self *Client) DeleteTask(namespace string, componentName string, taskName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, componentName, taskName)
	if podNames, err := kubernetes.GetPodNames(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		for _, podName := range podNames {
			if err := self.Exec(self.Namespace, podName, "operator", nil, nil, "rm", "--force", taskPath); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func (self *Client) ListTasks(namespace string, componentName string) ([]string, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		var buffer bytes.Buffer
		if err := self.Exec(self.Namespace, podName, "operator", nil, &buffer, "find", filepath.Join("/tasks", namespace, componentName), "-type", "f", "-printf", "%f\n"); err == nil {
			var names []string
			for _, filename := range strings.Split(strings.TrimRight(buffer.String(), "\n"), "\n") {
				names = append(names, filename)
			}
			return names, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) RunTask(namespace string, componentName string, taskName string, input ard.Value) (ard.Value, error) {
	if devices, err := self.ListDevices(); err == nil {
		task := make(ard.StringMap)
		task["name"] = taskName
		task["input"] = input

		component := make(ard.StringMap)
		component["name"] = componentName

		devices_ := make(ard.StringMap)
		for _, device := range devices.Items {
			device_ := resources.DeviceToARD(&device)
			delete(device_, "name")
			devices_[device.Name] = device_
		}

		args := make(ard.StringMap)
		args["task"] = task
		args["component"] = component
		args["devices"] = devices_

		if argsYaml, err := transcribe.EncodeYAML(args, " ", false); err == nil {
			// Default to same namespace as operator
			if namespace == "" {
				namespace = self.Namespace
			}

			appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
			taskPath := self.getTaskPath(namespace, componentName, taskName)
			if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
				var buffer bytes.Buffer
				if err := self.Exec(self.Namespace, podName, "operator", strings.NewReader(argsYaml), &buffer, taskPath); err == nil {
					value, _, err := ard.ReadYAML(&buffer, false)
					return value, err
				} else {
					if execError, ok := err.(*kubernetes.ExecError); ok {
						error_ := make(ard.StringMap)
						error_["message"] = execError.Err.Error()
						error_["stderr"] = execError.Stderr
						output := make(ard.StringMap)
						output["error"] = error_
						return output, nil
					} else {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) getTaskPath(namespace string, componentName string, taskName string) string {
	return filepath.Join("/tasks", namespace, componentName, taskName)
}
