package got

import (
	"context"
	"flag"
	"log"

	app_io "github.com/tennashi/got/io"
)

// Run is the entry point for the got command.
func Run(ctx context.Context, args []string, ioStream *app_io.Stream) int {
	g := newGot(ioStream)

	if err := g.run(ctx, args); err != nil {
		g.l.Println(err)
		return 1
	}
	return 0
}

type got struct {
	IOStream *app_io.Stream

	fs      *flag.FlagSet
	l       *log.Logger
	subCmds map[string]runner
}

func newGot(ioStream *app_io.Stream) *got {
	l := log.New(ioStream.Err, "[got] ", 0)
	fs := flag.NewFlagSet("got", flag.ContinueOnError)
	g := &got{fs: fs, l: l, IOStream: ioStream}

	subCmds := map[string]runner{
		"version": newVersionCmd(g),
		"help":    newHelpCmd(g),
		"clone":   newCloneCmd(g),
	}
	g.subCmds = subCmds
	return g
}

func (g *got) run(ctx context.Context, args []string) error {
	if err := g.fs.Parse(args); err != nil {
		return err
	}
	subCmd, ok := g.subCmds[g.fs.Arg(1)]
	if !ok {
		subCmd = g.subCmds["help"]
	}
	if err := subCmd.parse(args); err != nil {
		return err
	}
	return subCmd.run(ctx)
}

type runner interface {
	parse([]string) error
	run(context.Context) error
}
