package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

var configName = "config.toml"

type Config struct {
	paths    []string
	usedPath string

	Dotfiles Dotfiles
}

type Dotfiles struct {
	Local  string
	Remote string
}

func InitConfig(cfgFile string) (*Config, error) {
	config := &Config{}
	if cfgFile != "" {
		config.addPath(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
			config.addPath(filepath.Join(dir, "got", configName))
		}
		if dirs := os.Getenv("XDG_CONFIG_DIRS"); dirs != "" {
			for _, dir := range strings.Split(dirs, fmt.Sprintf("%c", filepath.ListSeparator)) {
				config.addPath(filepath.Join(dir, "got", configName))
			}
		}
		config.addPath(filepath.Join(home, ".config", "got", configName))
	}

	if err := config.load(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Write() error {
	if len(c.paths) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		path := filepath.Join(home, ".config", "got", configName)
		c.addPath(path)
		c.usedPath = path
	}

	f, err := os.OpenFile(c.usedPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) addPath(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	c.paths = append(c.paths, path)
}

func (c *Config) load() error {
	for _, path := range c.paths {
		if _, err := toml.DecodeFile(path, c); err != nil {
			return err
		}
		c.usedPath = path
		return nil
	}
	return errors.New("valid config is not exist")
}
