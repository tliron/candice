package commands

import (
	contextpkg "context"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
	"github.com/tliron/yamlkeys"
)

var inputs []string
var inputsUrl string

var inputValues = make(map[string]interface{})

func init() {
	taskCommand.AddCommand(taskRunCommand)
	taskRunCommand.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=YAML)")
	taskRunCommand.Flags().StringVarP(&inputsUrl, "inputs", "s", "", "load inputs from a PATH or URL to YAML content")
}

var taskRunCommand = &cobra.Command{
	Use:   "run [COMPONENT NAME] [TASK NAME]",
	Short: "Run a task for a component",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		RunTask(contextpkg.TODO(), args[0], args[1])
	},
}

func RunTask(context contextpkg.Context, componentName string, taskName string) {
	ParseInputs(context)
	result, err := NewClient().Candice().RunTask(namespace, componentName, taskName, inputValues)
	util.FailOnError(err)
	transcribe.Print(result, format, os.Stdout, strict, pretty)
}

func ParseInputs(context contextpkg.Context) {
	if inputsUrl != "" {
		log.Infof("load inputs from %q", inputsUrl)

		urlContext := exturl.NewContext()
		defer urlContext.Release()

		url, err := urlContext.NewValidURL(context, inputsUrl, nil)
		util.FailOnError(err)
		reader, err := url.Open(context)
		util.FailOnError(err)
		if closer, ok := reader.(io.Closer); ok {
			defer closer.Close()
		}
		data, err := yamlkeys.DecodeAll(reader)
		util.FailOnError(err)
		for _, data_ := range data {
			if map_, ok := data_.(ard.Map); ok {
				for key, value := range map_ {
					inputValues[yamlkeys.KeyString(key)] = value
				}
			} else {
				util.Failf("malformed inputs in %q", inputsUrl)
			}
		}
	}

	for _, input := range inputs {
		s := strings.SplitN(input, "=", 2)
		if len(s) != 2 {
			util.Failf("malformed input: %s", input)
		}
		value, _, err := ard.DecodeYAML(s[1], false)
		util.FailOnError(err)
		inputValues[s[0]] = value
	}
}
