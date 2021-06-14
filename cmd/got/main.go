package main

import (
	"os"

	"github.com/tennashi/got/cmd/got/internal/cli"
)

func main() {
	os.Exit(cli.Run())
}
