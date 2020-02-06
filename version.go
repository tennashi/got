package got

import (
	"fmt"
	"io"
	"runtime"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "none"
	goversion = runtime.Version()
)

func printVersion(outStream io.Writer) {
	fmt.Fprintf(outStream, "version: %v\n", version)
	fmt.Fprintf(outStream, "commit hash: %v\n", commit)
	fmt.Fprintf(outStream, "build date: %v\n", date)
	fmt.Fprintf(outStream, "go version: %v\n", goversion)
}
