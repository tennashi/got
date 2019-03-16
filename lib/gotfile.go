package lib

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Gotfile struct {
	Dotfile []Dotfile
}

type Dotfile struct {
	Dest string
	Src  string
}

func InitGotfile(dirPath string) (*Gotfile, error) {
	path := filepath.Join(dirPath, "Gotfile.toml")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	gotfile := &Gotfile{}
	viper.Unmarshal(gotfile)
	return gotfile, nil
}

func MakeGotfile(dirPath string) error {
	gotfilePath := filepath.Join(dirPath, "Gotfile.toml")
	if _, err := os.Stat(gotfilePath); err == nil {
		return errors.New("file exist")
	} else {
		if _, err := os.Create(gotfilePath); err != nil {
			return err
		}
	}
	return nil
}
