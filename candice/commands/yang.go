package commands

import (
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/cobra"

	//formatpkg "github.com/tliron/kutil/format"
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
		yang.AddPath("assets/yang/")
		entry, errs := yang.GetModule("ietf-interfaces@2017-12-16")
		if len(errs) > 0 {
			util.FailOnError(errs[0])
		}
		entry.Print(terminal.Stdout)
		//formatpkg.Print(entry.Name, format, terminal.Stdout, strict, pretty)
	},
}
