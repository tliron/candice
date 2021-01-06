package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	resources "github.com/tliron/candice/resources/candice.cloud/v1alpha1"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	deviceCommand.AddCommand(deviceListCommand)
}

var deviceListCommand = &cobra.Command{
	Use:   "list",
	Short: "List devices",
	Run: func(cmd *cobra.Command, args []string) {
		ListDevices()
	},
}

func ListDevices() {
	devices, err := NewClient().Candice().ListDevices()
	util.FailOnError(err)
	if len(devices.Items) == 0 {
		return
	}
	// TODO: sort devices by name? they seem already sorted!

	switch format {
	case "":
		table := terminal.NewTable(maxWidth, "Name", "Host", "Namespace", "Service", "Port", "LastError")
		for _, registry := range devices.Items {
			if registry.Spec.Direct != nil {
				table.Add(registry.Name, registry.Spec.Direct.Host, "", "", "", registry.Status.LastError)
			} else if registry.Spec.Indirect != nil {
				table.Add(registry.Name, "", registry.Spec.Indirect.Namespace, registry.Spec.Indirect.Service, fmt.Sprintf("%d", registry.Spec.Indirect.Port), registry.Status.LastError)
			}
		}
		table.Print()

	case "bare":
		for _, registry := range devices.Items {
			fmt.Fprintln(terminal.Stdout, registry.Name)
		}

	default:
		list := make(ard.List, len(devices.Items))
		for index, registry := range devices.Items {
			list[index] = resources.DeviceToARD(&registry)
		}
		formatpkg.Print(list, format, terminal.Stdout, strict, pretty)
	}
}
