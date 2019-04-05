package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	got "github.com/tennashi/got/lib"
)

// NewGetCmd returns get command.
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get package_name [package_manager]",
		Short: "install package from package repository",
		RunE:  doGet,
		Args:  cobra.RangeArgs(1, 2),
	}
	return cmd
}

func doGet(c *cobra.Command, args []string) error {
	var mgrName string
	if len(args) == 1 {
		if globalConfig == nil {
			return errors.New("config file unloaded")
		}
		mgrName = globalConfig.DefaultManager
	} else {
		mgrName = args[1]
	}
	manager := got.NewManager(mgrName)
	if manager == nil {
		return errors.New("unknown manager")
	}
	pkgName := args[0]
	err := manager.Install(pkgName)
	if err != nil {
		return err
	}
	dotfilesDir, err := got.ExpandPath(globalConfig.Dotfiles.Local)
	if err != nil {
		return err
	}
	gotfile, err := got.InitGotfile(dotfilesDir)
	if err != nil {
		return err
	}
	if err := gotfile.AddPackage(pkgName, mgrName); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(NewGetCmd())
}
