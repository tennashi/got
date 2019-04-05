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
	Package []Package
}

// Dotfile has your dotfile path.
type Dotfile struct {
	Dest string
	Src  string
}

// Package has the package name and the package manager command name.
type Package struct {
	Name    string
	Manager string
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

// AddPackage add the package config to the Gotfile.
func (g *Gotfile) AddPackage(name, manager string) error {
	pkgs := struct {
		Package []Package
	}{
		Package: []Package{
			{
				Name:    name,
				Manager: manager,
			},
		},
	}
	f, err := os.OpenFile(g.path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("# automatically added\n")
	if err := toml.NewEncoder(f).Encode(pkgs); err != nil {
		return err
	}
	f.WriteString("\n")
	return nil
}

// AddDotfile add the dotfile config to the Gotfile.
func (g *Gotfile) AddDotfile(src, dest string) error {
	dotfiles := struct {
		Dotfile []Dotfile
	}{
		Dotfile: []Dotfile{
			{
				Src:  src,
				Dest: dest,
			},
		},
	}
	f, err := os.OpenFile(g.path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("# automatically added\n")
	if err := toml.NewEncoder(f).Encode(dotfiles); err != nil {
		return err
	}
	f.WriteString("\n")
	return nil
}
