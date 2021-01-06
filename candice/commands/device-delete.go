package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	deviceCommand.AddCommand(deviceDeleteCommand)
	deviceDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all devices")
}

var deviceDeleteCommand = &cobra.Command{
	Use:   "delete [[DEVICE NAME]]",
	Short: "Delete devices",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			DeleteDevice(args[0])
		} else if all {
			DeleteAllRegistries()
		} else {
			util.Fail("must provide device name or specify \"--all\"")
		}
	},
}

func DeleteDevice(deviceName string) {
	// TODO: in cluster mode we must specify the namespace
	namespace := ""

	err := NewClient().Candice().DeleteDevice(namespace, deviceName)
	util.FailOnError(err)
}

func DeleteAllRegistries() {
	candice := NewClient().Candice()
	devices, err := candice.ListDevices()
	util.FailOnError(err)
	if len(devices.Items) > 0 {
		for _, registry := range devices.Items {
			log.Infof("deleting device: %s/%s", registry.Namespace, registry.Name)
			err := candice.DeleteDevice(registry.Namespace, registry.Name)
			util.FailOnError(err)
		}
	} else {
		log.Info("no devices to delete")
	}
}
