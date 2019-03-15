package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	dotfileConfig, err := setDotfileConfig(globalConfig, args)
	if err != nil {
		if err.Error() == "config file unloaded" {
			c.SetOutput(os.Stderr)
			c.Println("Config file unloaded.")
			c.Println("Please\n\t'got -c /path/to/config.toml sync'\nor\n\t'got sync https://your.remote.repository/dotfiles.git /path/to/dotfiles'.\n")
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

	gotFile, err := lib.InitGotFile(localDir)
	if err != nil {
		return err
	}

	for _, dotFile := range gotFile.DotFile {
		symLink, err := lib.NewSymLink(localDir, dotFile)
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

func setDotfileConfig(config *viper.Viper, args []string) (*dotfileRepo, error) {
	switch len(args) {
	case 0:
		if config.ConfigFileUsed() == "" {
			return nil, errors.New("config file unloaded")
		}
		localDir, err := lib.ExpandPath(config.GetString("dotfiles.local"))
		if err != nil {
			return nil, err
		}
		return &dotfileRepo{
			remoteURL: config.GetString("dotfiles.remote"),
			localDir:  localDir,
		}, nil
	case 1:
		if config.ConfigFileUsed() == "" {
			return nil, errors.New("config file unloaded")
		}
		localDir, err := lib.ExpandPath(config.GetString("dotfiles.local"))
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
