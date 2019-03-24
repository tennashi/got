package main

import (
	"runtime"

	"github.com/tennashi/got/cmd"
)

var (
	version   string = "dev"
	commit    string = "none"
	date      string = "none"
	goversion string = runtime.Version()
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Goversion = goversion

	cmd.Execute()
}
