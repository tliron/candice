package main

import (
	"github.com/tliron/candice/candice/commands"
	"github.com/tliron/kutil/util"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	commands.Execute()
	util.Exit(0)
}
