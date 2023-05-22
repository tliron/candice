package v1alpha1

import (
	"fmt"

	group "github.com/tliron/candice/resources/candice.puccini.cloud"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/kubernetes"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DeviceGVK = SchemeGroupVersion.WithKind(DeviceKind)

type DeviceProtocol string

const (
	DeviceKind     = "Device"
	DeviceListKind = "DeviceList"

	DeviceSingular  = "device"
	DevicePlural    = "devices"
	DeviceShortName = "dev"

	DeviceProtocolNETCONF  DeviceProtocol = "netconf"
	DeviceProtocolRESTCONF DeviceProtocol = "restconf"
)

//
// Device
//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Device struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceSpec   `json:"spec"`
	Status DeviceStatus `json:"status"`
}

type DeviceSpec struct {
	Protocol DeviceProtocol  `json:"protocol"`           // Device protocol
	Direct   *DeviceDirect   `json:"direct,omitempty"`   // Direct reference to device
	Indirect *DeviceIndirect `json:"indirect,omitempty"` // Indirect reference to device
}

type DeviceDirect struct {
	Host string `json:"host"` // Device host (either "host:port" or "host")
}

type DeviceIndirect struct {
	Namespace string `json:"namespace,omitempty"` // Namespace for service resource (optional; defaults to same namespace as this device)
	Service   string `json:"service"`             // Name of service resource
	Port      uint64 `json:"port"`                // Port to use with service
}

type DeviceStatus struct {
	LastError string `json:"lastError"` // Last error message on the device
}

//
// DeviceList
//

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeviceList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Device `json:"items"`
}

//
// DeviceCustomResourceDefinition
//

// See: assets/custom-resource-definitions.yaml

var DeviceResourcesName = fmt.Sprintf("%s.%s", DevicePlural, group.GroupName)

var DeviceCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: DeviceResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: DeviceSingular,
			Plural:   DevicePlural,
			Kind:     DeviceKind,
			ListKind: DeviceListKind,
			ShortNames: []string{
				DeviceShortName,
			},
			Categories: []string{
				"all", // will appear in "kubectl get all"
			},
		},
		Scope: apiextensions.NamespaceScoped,
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    Version,
				Served:  true,
				Storage: true, // one and only one version must be marked with storage=true
				Subresources: &apiextensions.CustomResourceSubresources{ // requires CustomResourceSubresources feature gate enabled
					Status: &apiextensions.CustomResourceSubresourceStatus{},
				},
				Schema: &apiextensions.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
						Description: "Candice device",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type:     "object",
								Required: []string{"protocol"},
								Properties: map[string]apiextensions.JSONSchemaProps{
									"protocol": {
										Description: "Device protocol",
										Type:        "string",
										Enum: []apiextensions.JSON{
											kubernetes.JSONString(DeviceProtocolNETCONF),
											kubernetes.JSONString(DeviceProtocolRESTCONF),
										},
									},
									"direct": {
										Description: "Direct reference to device",
										Type:        "object",
										Required:    []string{"host"},
										Properties: map[string]apiextensions.JSONSchemaProps{
											"host": {
												Description: "Device host (either \"host:port\" or \"host\")",
												Type:        "string",
											},
										},
									},
									"indirect": {
										Description: "Indirect reference to device",
										Type:        "object",
										Required:    []string{"service", "port"},
										Properties: map[string]apiextensions.JSONSchemaProps{
											"namespace": {
												Description: "Namespace for service resource (optional; defaults to same namespace as this device)",
												Type:        "string",
											},
											"service": {
												Description: "Name of service resource",
												Type:        "string",
											},
											"port": {
												Description: "Port to use with service",
												Type:        "integer",
											},
										},
									},
								},
								OneOf: []apiextensions.JSONSchemaProps{
									{
										Required: []string{"direct"},
									},
									{
										Required: []string{"indirect"},
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"lastError": {
										Description: "Last error message on the device",
										Type:        "string",
									},
								},
							},
						},
					},
				},
				AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
					{
						Name:     "Protocol",
						Type:     "string",
						JSONPath: ".spec.protocol",
					},
				},
			},
		},
	},
}

func DeviceToARD(device *Device) ard.StringMap {
	map_ := make(ard.StringMap)
	map_["name"] = device.Name
	map_["protocol"] = device.Spec.Protocol
	if device.Spec.Direct != nil {
		map_["direct"] = ard.StringMap{
			"host": device.Spec.Direct.Host,
		}
	} else if device.Spec.Indirect != nil {
		map_["indirect"] = ard.StringMap{
			"namespace": device.Spec.Indirect.Namespace,
			"service":   device.Spec.Indirect.Service,
			"port":      device.Spec.Indirect.Port,
		}
	}
	map_["lastError"] = device.Status.LastError
	return map_
}
