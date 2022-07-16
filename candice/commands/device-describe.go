package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	resources "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
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
		transcribe.Print(resources.DeviceToARD(device), format, terminal.Stdout, strict, pretty)
	} else {
		terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("Name"), terminal.Stylize.Value(device.Name))
		terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("Protocol"), terminal.Stylize.Value(string(device.Spec.Protocol)))

		if device.Spec.Direct != nil {
			terminal.Printf("  %s:\n", terminal.Stylize.TypeName("Direct"))
			if device.Spec.Direct.Host != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Host"), terminal.Stylize.Value(device.Spec.Direct.Host))
			}
		}

		if device.Spec.Indirect != nil {
			terminal.Printf("  %s:\n", terminal.Stylize.TypeName("Indirect"))
			if device.Spec.Indirect.Namespace != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Namespace"), terminal.Stylize.Value(device.Spec.Indirect.Namespace))
			}
			if device.Spec.Indirect.Service != "" {
				terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Service"), terminal.Stylize.Value(device.Spec.Indirect.Service))
			}
			terminal.Printf("    %s: %s\n", terminal.Stylize.TypeName("Port"), terminal.Stylize.Value(fmt.Sprintf("%d", device.Spec.Indirect.Port)))
		}

		if device.Status.LastError != "" {
			terminal.Printf("%s: %s\n", terminal.Stylize.TypeName("LastError"), terminal.Stylize.Value(device.Status.LastError))
		}
	}
}
