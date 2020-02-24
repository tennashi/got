package main

import (
	"context"
	"os"

	"github.com/tennashi/got"
)

func main() {
	os.Exit(got.Run(
		context.Background(),
		os.Args[1:],
		nil,
	))
}
