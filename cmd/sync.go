package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	got "github.com/tennashi/got/lib"
)

var writeConfig bool

// NewSyncCmd returns sync command.
func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [remote repository [local directory]]",
		Short: "clone your dotfiles repository",
		RunE:  doSync,
	}
	cmd.Flags().BoolVarP(&writeConfig, "write", "w", false, "overwrite config with command args")
	return cmd
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

	localDir, err := got.ExpandPath(dotfileConfig.localDir)
	if err != nil {
		return err
	}

	git := got.NewGit(dotfileConfig.remoteURL, localDir)
	if err := git.CloneOrPull(); err != nil {
		c.SetOutput(os.Stderr)
		c.Println("git:", err)
	}

	gotfile, err := got.InitGotfile(localDir)
	if err != nil {
		return err
	}

	pkgList := map[string][]string{}
	for _, pkg := range gotfile.Package {
		pkgList[pkg.Manager] = append(pkgList[pkg.Manager], pkg.Name)
	}
	for mgrName, pkgs := range pkgList {
		manager := got.NewManager(mgrName)
		if manager == nil {
			return errors.New("unknown manager")
		}
		err := manager.Install(pkgs...)
		if err != nil {
			return err
		}
	}

	for _, dotfile := range gotfile.Dotfile {
		symLink, err := got.NewSymLink(localDir, dotfile)
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

func setDotfileRepo(config *got.Config, args []string) (*dotfileRepo, error) {
	switch len(args) {
	case 0:
		if config == nil {
			return nil, errors.New("config file unloaded")
		}
		localDir, err := got.ExpandPath(config.Dotfiles.Local)
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
		if writeConfig {
			config.Dotfiles.Remote = args[0]
			config.Write()
		}

		localDir, err := got.ExpandPath(config.Dotfiles.Local)
		if err != nil {
			return nil, err
		}

		return &dotfileRepo{
			remoteURL: args[0],
			localDir:  localDir,
		}, nil
	case 2:
		if writeConfig {
			if config == nil {
				config = &got.Config{}
			}
			config.Dotfiles.Remote = args[0]
			config.Dotfiles.Local = args[1]
			config.Write()
		}

		localDir, err := got.ExpandPath(args[1])
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
