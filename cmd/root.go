package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
var globalConfig = viper.New()

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
	if cfgFile != "" {
		globalConfig.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			rootCmd.Println(err)
			os.Exit(1)
		}

		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS")

		globalConfig.SetConfigName("config")
		globalConfig.AddConfigPath(filepath.Join(xdgConfigHome, "got"))
		for _, dir := range strings.Split(xdgConfigDirs, fmt.Sprintf("%c", filepath.ListSeparator)) {
			globalConfig.AddConfigPath(filepath.Join(dir, "got"))
		}
		globalConfig.AddConfigPath(filepath.Join(home, ".config", "got"))
	}

	globalConfig.AutomaticEnv() // read in environment variables that match

	if err := globalConfig.ReadInConfig(); err == nil {
		rootCmd.Println("Using config file:", globalConfig.ConfigFileUsed())
	}
}
