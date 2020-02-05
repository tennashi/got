package got

import (
	"context"
	"fmt"
	"runtime"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "none"
	goversion = runtime.Version()
)

type versionCmd struct {
	rootCmd *got
}

func newVersionCmd(rootCmd *got) *versionCmd {
	return &versionCmd{rootCmd: rootCmd}
}

func (v *versionCmd) run(ctx context.Context, args []string) error {
	fmt.Fprintf(v.rootCmd.IOStream.Out, "version: %v\n", version)
	fmt.Fprintf(v.rootCmd.IOStream.Out, "commit hash: %v\n", commit)
	fmt.Fprintf(v.rootCmd.IOStream.Out, "build date: %v\n", date)
	fmt.Fprintf(v.rootCmd.IOStream.Out, "go version: %v\n", goversion)
	return nil
}
