package commands

import (
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/cobra"

	//"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

// EXPERIMENTAL

func init() {
	rootCommand.AddCommand(yangCommand)
}

var yangCommand = &cobra.Command{
	Use:   "yang",
	Short: "yang",
	Run: func(cmd *cobra.Command, args []string) {
		modules := yang.NewModules()
		modules.AddPath("assets/yang/")
		entry, errs := modules.GetModule("ietf-interfaces@2017-12-16")
		if len(errs) > 0 {
			util.FailOnError(errs[0])
		}
		entry.Print(terminal.Stdout)
		//transcribe.Print(entry.Name, format, terminal.Stdout, strict, pretty)
	},
}
