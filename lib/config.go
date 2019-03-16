package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Dotfiles Dotfiles
}

type Dotfiles struct {
	Local  string
	Remote string
}

func InitConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS")

		viper.SetConfigName("config")
		viper.AddConfigPath(filepath.Join(xdgConfigHome, "got"))
		for _, dir := range strings.Split(xdgConfigDirs, fmt.Sprintf("%c", filepath.ListSeparator)) {
			viper.AddConfigPath(filepath.Join(dir, "got"))
		}
		viper.AddConfigPath(filepath.Join(home, ".config", "got"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
