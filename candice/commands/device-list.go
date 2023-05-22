package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	resources "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
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
		table := terminal.NewTable(maxWidth, "Name", "Protocol", "Host", "Namespace", "Service", "Port", "LastError")
		for _, registry := range devices.Items {
			if registry.Spec.Direct != nil {
				table.Add(registry.Name, string(registry.Spec.Protocol), registry.Spec.Direct.Host, "", "", "", registry.Status.LastError)
			} else if registry.Spec.Indirect != nil {
				table.Add(registry.Name, string(registry.Spec.Protocol), "", registry.Spec.Indirect.Namespace, registry.Spec.Indirect.Service, fmt.Sprintf("%d", registry.Spec.Indirect.Port), registry.Status.LastError)
			}
		}
		table.Print()

	case "bare":
		for _, registry := range devices.Items {
			terminal.Println(registry.Name)
		}

	default:
		list := make(ard.List, len(devices.Items))
		for index, registry := range devices.Items {
			list[index] = resources.DeviceToARD(&registry)
		}
		transcribe.Print(list, format, os.Stdout, strict, pretty)
	}
}
