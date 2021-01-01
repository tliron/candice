package commands

import (
	contextpkg "context"

	"github.com/op/go-logging"
)

const toolName = "candice"

var context = contextpkg.TODO()

var log = logging.MustGetLogger(toolName)

var filePath string
var directoryPath string
var url string
var registry string
var tail int
var follow bool
var all bool
var sourceRegistry string
var wait bool
