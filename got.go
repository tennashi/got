package got

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
)

// Run is the entry point for the got command.
func Run(ctx context.Context, args []string, ioStream *ioStream) int {
	if ioStream == nil {
		ioStream = defaultIOStream
	}
	g := newGot(ioStream)

	if err := g.parse(args); err != nil {
		g.l.Println(err)
		return 1
	}
	if err := g.run(ctx); err != nil {
		g.l.Println(err)
		return 1
	}
	return 0
}

type got struct {
	IOStream *ioStream

	fs      *flag.FlagSet
	l       *log.Logger
	subCmds map[string]runner

	curSubCmd runner
	err       error
}

func newGot(ioStream *ioStream) *got {
	l := log.New(ioStream.Err, "[got] ", 0)

	fs := flag.NewFlagSet("got", flag.ContinueOnError)
	fs.SetOutput(ioStream.Err)
	fs.Usage = func() {
		fmt.Fprint(ioStream.Out, help)
		fs.PrintDefaults()
	}

	g := &got{fs: fs, l: l, IOStream: ioStream}

	subCmds := map[string]runner{
		"version": newVersionCmd(g),
		"help":    newHelpCmd(g),
		"get":     newGetCmd(g),
		"remove":  newRemoveCmd(g),
	}
	g.subCmds = subCmds
	return g
}

func (g *got) parse(args []string) error {
	if err := g.fs.Parse(args); err != nil {
		return err
	}
	if g.fs.NArg() == 0 {
		g.err = errors.New("must specify sub command")
		return g.showHelp()
	}

	var ok bool
	g.curSubCmd, ok = g.subCmds[g.fs.Arg(0)]
	if !ok {
		g.err = fmt.Errorf("no such command: %s", g.fs.Arg(0))
		return g.showHelp()
	}
	if err := g.curSubCmd.parse(g.fs.Args()[1:]); err != nil {
		g.err = err
		return g.showHelp()
	}
	return nil
}

func (g *got) run(ctx context.Context) error {
	return g.curSubCmd.run(ctx)
}

type runner interface {
	parse([]string) error
	run(context.Context) error
}

func (g *got) showHelp() error {
	g.fs.Usage()
	if g.err != nil {
		return g.err
	}
	return nil
}

var help = `got - go packages manager

Usage:
  got command [arguments]

Commands:
  help
    print this help

  version
    print got command version

  get [-lu] [-c command] package
    install the package

  remove [package|command]
    remove the package

`
