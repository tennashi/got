package got

import (
	"context"
	"errors"
	"flag"
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

func newHelpCmd(rootCmd *got) *helpCmd {
	fs := flag.NewFlagSet("got-help", flag.ContinueOnError)
	return &helpCmd{rootCmd: rootCmd, fs: fs}
}

func (c *helpCmd) parse(args []string) error {
	return c.fs.Parse(args)
}

func (c *helpCmd) run(ctx context.Context) error {
	return c.rootCmd.showHelp()
}

type getCmd struct {
	rootCmd *got
	fs      *flag.FlagSet

	isUpdate bool
	isList   bool
	cmdName  string
	pkgName  string
}

func newGetCmd(rootCmd *got) *getCmd {
	fs := flag.NewFlagSet("got-get", flag.ContinueOnError)
	fs.SetOutput(rootCmd.IOStream.Err)

	cmd := getCmd{
		rootCmd: rootCmd,
		fs:      fs,
	}

	fs.BoolVar(&cmd.isUpdate, "u", false, "update")
	fs.BoolVar(&cmd.isList, "l", false, "list")
	fs.StringVar(&cmd.cmdName, "c", "", "command name")
	return &cmd
}

func (c *getCmd) parse(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if c.fs.NArg() == 0 && !c.isList {
		return errors.New("must specify package name")
	}
	c.pkgName = c.fs.Arg(0)
	return nil
}

func (c *getCmd) run(ctx context.Context) error {
	if c.isList {
		return getAll(c.rootCmd.IOStream, c.isUpdate)
	}

	return get(c.rootCmd.IOStream, c.pkgName, c.cmdName, c.isUpdate)
}

type removeCmd struct {
	rootCmd *got
	fs      *flag.FlagSet

	targetName string
}

func newRemoveCmd(rootCmd *got) *removeCmd {
	fs := flag.NewFlagSet("got-delete", flag.ContinueOnError)
	fs.SetOutput(rootCmd.IOStream.Err)

	return &removeCmd{
		rootCmd: rootCmd,
		fs:      fs,
	}
}

func (c *removeCmd) parse(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if c.fs.NArg() == 0 {
		return errors.New("must specify package/command name")
	}
	c.targetName = c.fs.Arg(0)
	return nil
}

func (c *removeCmd) run(ctx context.Context) error {
	return remove(c.rootCmd.IOStream, c.targetName)
}
