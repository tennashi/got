package main

import "github.com/tennashi/got/cmd"

var (
	version   string
	hash      string
	builddate string
	goversion string
)

func main() {
	cmd.AppVer = version
	cmd.Hash = hash
	cmd.Builddate = builddate
	cmd.Goversion = goversion

	cmd.Execute()
}
