package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tennashi/got/lib"
)

// AppVer is application version
// Hash is commit hash
// Builddate is date of build
// Goversion is go version
var (
	AppVer    string
	Hash      string
	Builddate string
	Goversion string
)

var cfgFile string
var globalConfig *lib.Config

var rootCmd = &cobra.Command{
	Use:   "got [command]",
	Short: "dotfile manager written in Go",
	Long: `got is dotfile manager.
This application is a tool to manage your (local or remote) dotfiles repository.`,
}

// Execute is entry point
func Execute() {
	rootCmd.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.got/config.yaml)")
}

func initConfig() {
	var err error
	globalConfig, err = lib.InitConfig(cfgFile)
	if err != nil {
		os.Exit(1)
	}
}
