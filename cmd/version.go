package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run:   doVersion,
}

func doVersion(c *cobra.Command, Args []string) {
	fmt.Println("version:", AppVer)
	fmt.Println("commit hash:", Hash)
	fmt.Println("build date:", Builddate)
	fmt.Println("go version:", Goversion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
