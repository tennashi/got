package lib

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var gotfileName = "Gotfile.toml"

type Gotfile struct {
	path    string
	Dotfile []Dotfile
}

type Dotfile struct {
	Dest string
	Src  string
}

func InitGotfile(dirPath string) (*Gotfile, error) {
	gotfile := &Gotfile{}
	path := filepath.Join(dirPath, gotfileName)
	gotfile.setPath(path)
	if err := gotfile.load(); err != nil {
		return nil, err
	}

	return gotfile, nil
}

func MakeGotfile(dirPath string) error {
	path := filepath.Join(dirPath, gotfileName)
	if _, err := os.Stat(path); err == nil {
		return errors.New("file exist")
	} else {
		if _, err := os.Create(path); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gotfile) setPath(path string) {
	g.path = path
}

func (g *Gotfile) load() error {
	if _, err := toml.DecodeFile(g.path, g); err != nil {
		return err
	}
	return nil
}
