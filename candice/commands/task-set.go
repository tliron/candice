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
	Use:   "set [DEVICE NAME] [TASK NAME]",
	Short: "Sets a task for a device",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		SetTask(args[0], args[1])
	},
}

func SetTask(deviceName string, taskName string) {
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

	err = NewClient().Client().SetTask(namespace, deviceName, taskName, util.BytesToString(bytes))
	util.FailOnError(err)
}
