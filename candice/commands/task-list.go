package commands

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

func init() {
	taskCommand.AddCommand(taskListCommand)
}

var taskListCommand = &cobra.Command{
	Use:   "list [DEVICE NAME]",
	Short: "List tasks for a device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ListTasks(args[0])
	},
}

func ListTasks(deviceName string) {
	tasks, err := NewClient().Client().ListTasks(namespace, deviceName)
	util.FailOnError(err)
	if len(tasks) == 0 {
		return
	}
	sort.Strings(tasks)

	switch format {
	case "":
		// TODO fill table
		table := terminal.NewTable(maxWidth, "Name", "Server", "Namespace")
		for _, task := range tasks {
			table.Add(task, "TODO", "TODO")
		}
		table.Print()

	case "bare":
		for _, task := range tasks {
			fmt.Fprintln(terminal.Stdout, task)
		}

	default:
		list := make(ard.List, len(tasks))
		for index, task := range tasks {
			map_ := make(ard.StringMap)
			map_["Name"] = task
			map_["Server"] = ""
			map_["Namespace"] = ""
			list[index] = map_
		}
		formatpkg.Print(list, format, terminal.Stdout, strict, pretty)
	}
}
