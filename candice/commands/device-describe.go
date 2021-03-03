package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	resources "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	deviceCommand.AddCommand(deviceDescribeCommand)
}

var deviceDescribeCommand = &cobra.Command{
	Use:   "describe [DEVICE NAME]",
	Short: "Describe a device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DescribeDevice(args[0])
	},
}

func DescribeDevice(deviceName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	device, err := NewClient().Candice().GetDevice(namespace, deviceName)
	util.FailOnError(err)

	if format != "" {
		formatpkg.Print(resources.DeviceToARD(device), format, terminal.Stdout, strict, pretty)
	} else {
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("Name"), terminal.StyleValue(device.Name))
		fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("Protocol"), terminal.StyleValue(string(device.Spec.Protocol)))

		if device.Spec.Direct != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.StyleTypeName("Direct"))
			if device.Spec.Direct.Host != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Host"), terminal.StyleValue(device.Spec.Direct.Host))
			}
		}

		if device.Spec.Indirect != nil {
			fmt.Fprintf(terminal.Stdout, "  %s:\n", terminal.StyleTypeName("Indirect"))
			if device.Spec.Indirect.Namespace != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Namespace"), terminal.StyleValue(device.Spec.Indirect.Namespace))
			}
			if device.Spec.Indirect.Service != "" {
				fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Service"), terminal.StyleValue(device.Spec.Indirect.Service))
			}
			fmt.Fprintf(terminal.Stdout, "    %s: %s\n", terminal.StyleTypeName("Port"), terminal.StyleValue(fmt.Sprintf("%d", device.Spec.Indirect.Port)))
		}

		if device.Status.LastError != "" {
			fmt.Fprintf(terminal.Stdout, "%s: %s\n", terminal.StyleTypeName("LastError"), terminal.StyleValue(device.Status.LastError))
		}
	}
}
