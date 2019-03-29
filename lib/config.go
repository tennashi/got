package got

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

// Config represent config file.
type Config struct {
	paths    []string
	usedPath string

	Dotfiles Dotfiles
}

// Dotfiles have the dotfile repository location.
type Dotfiles struct {
	Local  string
	Remote string
}

// InitConfig initialize Config type.
// This function tries to read settings in the following order:
//   * The file specified by `-c` option.
//   * $XDG_CONFIG_HOME/got/config.toml
//   * $xdg_config_dir/got/config.toml ($xdg_config_dir is the path contained in $XDG_CONFIG_DIRS)
//   * ~/.config/got/config.toml
func InitConfig(cfgFile string) (*Config, error) {
	config := &Config{}

	if err := config.setPaths(cfgFile); err != nil {
		return nil, err
	}

	if err := config.load(); err != nil {
		return nil, err
	}

	return config, nil
}

// Write writes the current config to the currently used config file.
// If command haven't read the config file, this method write the config to `~/.config/got/config.toml`.
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

func (c *Config) setPaths(cfgFile string) error {
	if cfgFile != "" {
		c.addPath(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
			c.addPath(filepath.Join(dir, "got", configName))
		}
		if dirs := os.Getenv("XDG_CONFIG_DIRS"); dirs != "" {
			for _, dir := range strings.Split(dirs, fmt.Sprintf("%c", filepath.ListSeparator)) {
				c.addPath(filepath.Join(dir, "got", configName))
			}
		}
		c.addPath(filepath.Join(home, ".config", "got", configName))
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
