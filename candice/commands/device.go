package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(deviceCommand)
}

var deviceCommand = &cobra.Command{
	Use:   "device",
	Short: "Work with devices",
}
