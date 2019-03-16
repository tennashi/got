package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tennashi/got/lib"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialize your dotfiles repository",
		RunE:  doInit,
	}
}

func doInit(c *cobra.Command, args []string) error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := lib.MakeGotfile(curDir); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(NewInitCmd())
}
