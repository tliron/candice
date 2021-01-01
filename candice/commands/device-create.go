package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

var host string
var serviceNamespace string
var service string
var port uint64

func init() {
	deviceCommand.AddCommand(deviceCreateCommand)
	deviceCreateCommand.Flags().StringVarP(&host, "host", "", "", "device host (\"host\" or \"host:port\")")
	deviceCreateCommand.Flags().StringVarP(&serviceNamespace, "service-namespace", "", "", "device service namespace name (defaults to device namespace)")
	deviceCreateCommand.Flags().StringVarP(&service, "service", "", "", "device service name")
	deviceCreateCommand.Flags().Uint64VarP(&port, "port", "", 5000, "device service port")
}

var deviceCreateCommand = &cobra.Command{
	Use:   "create [DEVICE NAME]",
	Short: "Create a device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CreateDevice(args[0])
	},
}

func CreateDevice(deviceName string) {
	if host != "" {
		if service != "" {
			failDeviceCreate()
		}
	} else if service != "" {
		if host != "" {
			failDeviceCreate()
		}
	} else {
		failDeviceCreate()
	}

	client := NewClient().Client()
	var err error
	if service != "" {
		_, err = client.CreateDeviceIndirect(namespace, deviceName, serviceNamespace, service, port)
	} else {
		_, err = client.CreateDeviceDirect(namespace, deviceName, host)
	}
	util.FailOnError(err)
}

func failDeviceCreate() {
	util.Fail("must specify only one of \"--host\" or \"--service\"")
}
