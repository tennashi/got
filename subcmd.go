package got

import (
	"context"
	"flag"
	"fmt"
)

type versionCmd struct {
	rootCmd *got
}

func newVersionCmd(rootCmd *got) *versionCmd {
	return &versionCmd{rootCmd: rootCmd}
}

func (v *versionCmd) parse(args []string) error {
	return nil
}

func (v *versionCmd) run(ctx context.Context) error {
	printVersion(v.rootCmd.IOStream.Out)
	return nil
}

type helpCmd struct {
	rootCmd *got
	fs      *flag.FlagSet
}

var help = `got - packages manager

Usage:
  got command [arguments]

Commands:
  version
    print got command version

  get [manager] package...
    install the package using the manager

  update [manager] [package]
    update packages using the manager

  remove [manager] [package]
    remove the package using the manager

`

func newHelpCmd(rootCmd *got) *helpCmd {
	fs := flag.NewFlagSet("got-help", flag.ContinueOnError)
	fs.SetOutput(rootCmd.IOStream.Err)
	fs.Usage = func() {
		fmt.Fprint(rootCmd.IOStream.Out, help)
		fs.PrintDefaults()
	}
	return &helpCmd{rootCmd: rootCmd, fs: fs}
}

func (c *helpCmd) parse(args []string) error {
	return c.fs.Parse(args)
}

func (c *helpCmd) run(ctx context.Context) error {
	c.fs.Usage()
	if c.fs.NArg() < 2 {
		return fmt.Errorf("must specify sub-command")
	}
	subCmd := c.fs.Arg(1)
	if subCmd != "help" {
		return fmt.Errorf("no such command: %s", subCmd)
	}
	return nil
}
