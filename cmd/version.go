package cmd

import (
	"github.com/spf13/cobra"
)

// NewVersionCmd returns version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run:   doVersion,
	}
}

func doVersion(c *cobra.Command, args []string) {
	c.Println("version:", Version)
	c.Println("commit hash:", Commit)
	c.Println("build date:", Date)
	c.Println("go version:", Goversion)
}

func init() {
	rootCmd.AddCommand(NewVersionCmd())
}
