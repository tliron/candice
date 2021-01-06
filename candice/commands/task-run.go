package commands

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
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
		RunTask(args[0], args[1])
	},
}

func RunTask(componentName string, taskName string) {
	ParseInputs()
	result, err := NewClient().Candice().RunTask(namespace, componentName, taskName, inputValues)
	util.FailOnError(err)
	result, _ = ard.ToStringMaps(result)
	formatpkg.Print(result, format, terminal.Stdout, strict, pretty)
}

func ParseInputs() {
	if inputsUrl != "" {
		log.Infof("load inputs from %q", inputsUrl)

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		url, err := urlpkg.NewValidURL(inputsUrl, nil, urlContext)
		util.FailOnError(err)
		reader, err := url.Open()
		util.FailOnError(err)
		if closer, ok := reader.(io.Closer); ok {
			defer closer.Close()
		}
		data, err := formatpkg.ReadAllYAML(reader)
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
		value, err := formatpkg.DecodeYAML(s[1])
		util.FailOnError(err)
		inputValues[s[0]] = value
	}
}
