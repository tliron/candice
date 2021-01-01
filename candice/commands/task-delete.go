package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	taskCommand.AddCommand(taskDeleteCommand)
	taskDeleteCommand.Flags().BoolVarP(&all, "all", "a", false, "delete all tasks")
}

var taskDeleteCommand = &cobra.Command{
	Use:   "delete [DEVICE NAME] [[TASK NAME]]",
	Short: "Delete a task for a device",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			DeleteTask(args[0], args[1])
		} else if all {
			DeleteAllTasks(args[0])
		} else {
			util.Fail("must provide task name or specify \"--all\"")
		}
	},
}

func DeleteTask(deviceName string, taskName string) {
	err := NewClient().Client().DeleteTask(namespace, deviceName, taskName)
	util.FailOnError(err)
}

func DeleteAllTasks(deviceName string) {
	client := NewClient().Client()
	tasks, err := client.ListTasks(namespace, deviceName)
	util.FailOnError(err)
	for _, task := range tasks {
		log.Infof("deleting task: %s", task)
		err := client.DeleteTask(namespace, deviceName, task)
		util.FailOnError(err)
	}
}
