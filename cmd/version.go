package cmd

import (
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run:   doVersion,
	}
}

func doVersion(c *cobra.Command, args []string) {
	c.Println("version:", AppVer)
	c.Println("commit hash:", Hash)
	c.Println("build date:", Builddate)
	c.Println("go version:", Goversion)
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}
