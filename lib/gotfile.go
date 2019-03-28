package got

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var gotfileName = "Gotfile.toml"

// Gotfile represent Gotfile.
type Gotfile struct {
	path    string
	Dotfile []Dotfile
}

// Dotfile have your dotfile path.
type Dotfile struct {
	Dest string
	Src  string
}

// InitGotfile initialize Gotfile type from the dotfiles directory path.
func InitGotfile(dirPath string) (*Gotfile, error) {
	gotfile := &Gotfile{}
	path := filepath.Join(dirPath, gotfileName)
	gotfile.setPath(path)
	if err := gotfile.load(); err != nil {
		return nil, err
	}

	return gotfile, nil
}

// MakeGotfile creates Gotfile.toml to the dirPath.
func MakeGotfile(dirPath string) error {
	path := filepath.Join(dirPath, gotfileName)
	if _, err := os.Stat(path); err == nil {
		return errors.New("file exist")
	}
	if _, err := os.Create(path); err != nil {
		return err
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
