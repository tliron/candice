package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/kubernetes"
)

var permissions int64 = 0700

func (self *Client) GetTask(namespace string, deviceName string, taskName string) (string, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, deviceName, taskName)
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

func (self *Client) SetTask(namespace string, deviceName string, taskName string, content string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, deviceName, taskName)
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

func (self *Client) DeleteTask(namespace string, deviceName string, taskName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	taskPath := self.getTaskPath(namespace, deviceName, taskName)
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

func (self *Client) ListTasks(namespace string, deviceName string) ([]string, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
	if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
		var buffer bytes.Buffer
		if err := self.Exec(self.Namespace, podName, "operator", nil, &buffer, "find", filepath.Join("/tasks", namespace, deviceName), "-type", "f", "-printf", "%f\n"); err == nil {
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

func (self *Client) RunTask(namespace string, deviceName string, taskName string, input ard.Value) (ard.Value, error) {
	if inputYaml, err := format.EncodeYAML(input, " ", false); err == nil {
		// Default to same namespace as operator
		if namespace == "" {
			namespace = self.Namespace
		}

		appName := fmt.Sprintf("%s-%s", self.NamePrefix, "operator")
		taskPath := self.getTaskPath(namespace, deviceName, taskName)
		if podName, err := kubernetes.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, appName); err == nil {
			var buffer bytes.Buffer
			if err := self.Exec(self.Namespace, podName, "operator", strings.NewReader(inputYaml), &buffer, taskPath); err == nil {
				return format.DecodeYAML(buffer.String())
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

func (self *Client) getTaskPath(namespace string, deviceName string, taskName string) string {
	return filepath.Join("/tasks", namespace, deviceName, taskName)
}
