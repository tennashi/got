package main

import (
	"context"
	"os"

	"github.com/tennashi/got"
	app_io "github.com/tennashi/got/io"
)

func main() {
	ioStream := &app_io.Stream{
		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
	}
	os.Exit(got.Run(
		context.Background(),
		os.Args,
		ioStream,
	))

}
