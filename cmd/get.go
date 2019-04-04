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
	var manager got.Manager
	if len(args) == 1 {
		if globalConfig == nil {
			return errors.New("config file unloaded")
		}
		manager = got.NewManager(globalConfig.DefaultManager)
	} else {
		manager = got.NewManager(args[1])
	}
	if manager == nil {
		return errors.New("unknown manager")
	}
	err := manager.Install(args[0])
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(NewGetCmd())
}
