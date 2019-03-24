package main

import (
	"runtime"

	"github.com/tennashi/got/cmd"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "none"
	goversion = runtime.Version()
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Goversion = goversion

	cmd.Execute()
}
