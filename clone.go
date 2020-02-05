package got

import (
	"context"
	"errors"
	"flag"

	"github.com/tennashi/got/repository"
)

type cloneCmd struct {
	rootCmd *got

	fs *flag.FlagSet
}

func newCloneCmd(rootCmd *got) *cloneCmd {
	fs := flag.NewFlagSet("got-clone", flag.ContinueOnError)
	return &cloneCmd{rootCmd: rootCmd, fs: fs}
}

func (c *cloneCmd) run(ctx context.Context, args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	if c.fs.NArg() < 3 {
		return errors.New("must specify the dotfiles repository URL")
	}
	repoName := args[2]
	return c.clone(repoName)
}

func (c *cloneCmd) clone(repoName string) error {
	repo, err := repository.New(c.rootCmd.IOStream, repoName)
	if err != nil {
		return err
	}
	return repo.Clone()
}
