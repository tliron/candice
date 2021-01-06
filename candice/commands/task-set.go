package commands

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
)

func init() {
	taskCommand.AddCommand(taskSetCommand)
	taskSetCommand.Flags().StringVarP(&filePath, "file", "f", "", "path to a local task file (will be uploaded)")
}

var taskSetCommand = &cobra.Command{
	Use:   "set [COMPONENT NAME] [TASK NAME]",
	Short: "Sets a task for a component",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetTask(args[0], args[1])
	},
}

func SetTask(componentName string, taskName string) {
	var reader io.Reader
	if filePath != "" {
		var err error
		reader, err = os.Open(filePath)
		util.FailOnError(err)
	} else {
		reader = os.Stdin
	}

	bytes, err := ioutil.ReadAll(reader)
	util.FailOnError(err)

	err = NewClient().Candice().SetTask(namespace, componentName, taskName, util.BytesToString(bytes))
	util.FailOnError(err)
}
