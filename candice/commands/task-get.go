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
	Use:   "get [COMPONENT NAME] [TASK NAME]",
	Short: "Gets a task for a component",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		GetTask(args[0], args[1])
	},
}

func GetTask(componentName string, taskName string) {
	task, err := NewClient().Candice().GetTask(namespace, componentName, taskName)
	util.FailOnError(err)
	fmt.Fprintln(terminal.Stdout, task)
}
