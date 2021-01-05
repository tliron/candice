package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	operatorCommand.AddCommand(operatorLogsCommand)
	operatorLogsCommand.Flags().IntVarP(&tail, "tail", "t", -1, "number of most recent lines to print (<0 means all lines)")
	operatorLogsCommand.Flags().BoolVarP(&follow, "follow", "f", false, "keep printing incoming logs")
}

var operatorLogsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of the Candice operator",
	Run: func(cmd *cobra.Command, args []string) {
		Logs("operator", "operator")
	},
}
