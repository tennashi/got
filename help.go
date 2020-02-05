package got

import (
	"context"
	"flag"
	"fmt"
)

var help = `got - dotfiles and packages manager

Usage:
  got command [arguments]

Commands:
  version
    print got command version

  clone [repo_name]
    clone your dotfiles repository

  push
    push dotfiles update to your dotfiles repository

  get [manager] package...
    install the package using the manager

  update [manager] [package]
    update packages using the manager

  remove [manager] [package]
    remove the package using the manager

  link
    create symbolic links

  clean
    remove all symbolic links

`

type helpCmd struct {
	rootCmd *got
	fs      *flag.FlagSet
}

func newHelpCmd(rootCmd *got) *helpCmd {
	fs := flag.NewFlagSet("got-help", flag.ContinueOnError)
	fs.SetOutput(rootCmd.IOStream.Err)
	fs.Usage = func() {
		fmt.Fprint(rootCmd.IOStream.Out, help)
	}
	return &helpCmd{rootCmd: rootCmd, fs: fs}
}

func (c helpCmd) run(ctx context.Context, args []string) error {
	c.fs.Usage()
	if len(args) < 2 {
		return fmt.Errorf("must specify sub-command")
	}
	if args[1] != "help" {
		return fmt.Errorf("no such command: %s", args[1])
	}
	return nil
}
