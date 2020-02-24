package got

import (
	"fmt"
	"io"
	"runtime"
)

var (
	version   = "v0.0.1"
	goversion = runtime.Version()
)

func printVersion(outStream io.Writer) {
	fmt.Fprintf(outStream, "version: %v\n", version)
	fmt.Fprintf(outStream, "go version: %v\n", goversion)
}
