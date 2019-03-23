package main

import "github.com/tennashi/got/cmd"

var (
	version   string
	commit    string
	date      string
	goversion string
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Goversion = goversion

	cmd.Execute()
}
