package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	taskCommand.AddCommand(taskGetCommand)
}

var taskGetCommand = &cobra.Command{
	Use:   "get [DEVICE NAME] [TASK NAME]",
	Short: "Gets a task for a device",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		GetTask(args[0], args[1])
	},
}

func GetTask(deviceName string, taskName string) {
	task, err := NewClient().Client().GetTask(namespace, deviceName, taskName)
	util.FailOnError(err)
	fmt.Fprintln(terminal.Stdout, task)
}
