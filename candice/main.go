package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/candice/candice/commands"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}
