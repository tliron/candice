package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(taskCommand)
}

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Work with tasks",
}
