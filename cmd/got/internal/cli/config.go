package cli

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
)

func defaultConfigPath() (string, error) {
	return xdg.ConfigFile("got/config.toml")
}

func defaultBinDir() (string, error) {
	gobin := os.Getenv("GOBIN")
	if gobin != "" {
		return gobin, nil
	}

	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return filepath.Join(gopath, "bin"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "go/bin"), nil
}

func defaultDataDir() string {
	return filepath.Join(xdg.DataHome, "got")
}

type ConfigFile struct {
	InstallAll bool   `toml:"install_all"`
	BinDir     string `toml:"bin_dir"`
	DataDir    string `toml:"data_dir"`
}

func NewDefaultConfigFile() (*ConfigFile, error) {
	defaultBinDir, err := defaultBinDir()
	if err != nil {
		return nil, err
	}

	return &ConfigFile{
		InstallAll: true,
		BinDir:     defaultBinDir,
		DataDir:    defaultDataDir(),
	}, nil
}

func NewConfigFile(path string) (*ConfigFile, error) {
	if path == "" {
		var err error
		path, err = defaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			defaultConfig, derr := NewDefaultConfigFile()
			if derr != nil {
				return nil, derr
			}
			return defaultConfig, err
		}
		return nil, err
	}

	configFile := ConfigFile{}
	if err := toml.NewDecoder(f).Decode(&configFile); err != nil {
		return nil, err
	}

	return &configFile, nil
}
