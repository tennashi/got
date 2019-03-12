package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialize your dotfiles repository",
		Run:   doInit,
	}
}

func doInit(c *cobra.Command, args []string) {
	curDir, err := os.Getwd()
	if err != nil {
		c.Println(err)
	}

	gotFilePath := filepath.Join(curDir, "Gotfile.toml")
	if _, err := os.Stat(gotFilePath); err == nil {
		c.Println("got:", gotFilePath, "exist")
	} else {
		if _, err := os.Create(gotFilePath); err != nil {
			c.Println(err)
		}
		c.Println("got:", gotFilePath, "created")
	}
}

func init() {
	rootCmd.AddCommand(NewInitCmd())
}
