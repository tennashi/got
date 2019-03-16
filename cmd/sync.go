package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/tennashi/got/lib"
)

func NewSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync [remote repository [local directory]]",
		Short: "clone your dotfiles repository",
		RunE:  doSync,
	}
}

func doSync(c *cobra.Command, args []string) error {
	dotfileConfig, err := setDotfileRepo(globalConfig, args)
	if err != nil {
		if err.Error() == "config file unloaded" {
			c.SetOutput(os.Stderr)
			c.Println("Config file unloaded.")
			c.Println("Please\n\t'got -c /path/to/config.toml sync'\nor\n\t'got sync https://your.remote.repository/dotfiles.git /path/to/dotfiles'.")
			c.Usage()
		}
		return err
	}

	localDir, err := lib.ExpandPath(dotfileConfig.localDir)
	if err != nil {
		return err
	}

	git := lib.NewGit(dotfileConfig.remoteURL, localDir)
	if err := git.Clone(); err != nil {
		c.SetOutput(os.Stderr)
		c.Println("git:", err)
	}

	gotfile, err := lib.InitGotfile(localDir)
	if err != nil {
		return err
	}

	for _, dotfile := range gotfile.Dotfile {
		symLink, err := lib.NewSymLink(localDir, dotfile)
		if err != nil {
			c.SetOutput(os.Stderr)
			c.Println("symlink:", err)
			continue
		}
		if err := symLink.Make(); err != nil {
			c.SetOutput(os.Stderr)
			c.Println("symlink:", err)
		}
	}
	return nil
}

type dotfileRepo struct {
	remoteURL string
	localDir  string
}

func setDotfileRepo(config *lib.Config, args []string) (*dotfileRepo, error) {
	switch len(args) {
	case 0:
		if config == nil {
			return nil, errors.New("config file unloaded")
		}
		localDir, err := lib.ExpandPath(config.Dotfiles.Local)
		if err != nil {
			return nil, err
		}
		return &dotfileRepo{
			remoteURL: config.Dotfiles.Remote,
			localDir:  localDir,
		}, nil
	case 1:
		if config == nil {
			return nil, errors.New("config file unloaded")
		}
		localDir, err := lib.ExpandPath(config.Dotfiles.Local)
		if err != nil {
			return nil, err
		}
		return &dotfileRepo{
			remoteURL: args[0],
			localDir:  localDir,
		}, nil
	case 2:
		localDir, err := lib.ExpandPath(args[1])
		if err != nil {
			return nil, err
		}
		return &dotfileRepo{
			remoteURL: args[0],
			localDir:  localDir,
		}, nil
	default:
		return nil, errors.New("invalid argument")
	}
}

func init() {
	rootCmd.AddCommand(NewSyncCmd())
}
