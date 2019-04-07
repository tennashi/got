package got

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var gotfileName = "Gotfile.toml"

// Gotfile represent Gotfile.
type Gotfile struct {
	path    string
	Dotfile []*Dotfile          `toml:"dotfile"`
	Package map[string]*Package `toml:"package"`
}

// Dotfile has your dotfile path.
type Dotfile struct {
	Dest string `toml:"dest"`
	Src  string `toml:"src"`
}

// Package has the package name and the package manager command name.
type Package struct {
	Names []string `toml:"names"`
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

func (g *Gotfile) save() error {
	if err := os.Remove(g.path); err != nil {
		return err
	}
	f, err := os.OpenFile(g.path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("hoge")
		return err
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(g); err != nil {
		return err
	}
	return nil
}

// AddPackage add the package config to the Gotfile.
func (g *Gotfile) AddPackage(name, manager string) error {
	if _, ok := g.Package[manager]; !ok {
		g.Package = map[string]*Package{
			manager: &Package{Names: []string{name}},
		}
	} else {
		tmp := append(g.Package[manager].Names, name)
		m := make(map[string]struct{})
		unique := make([]string, 0, len(tmp))
		for _, name := range tmp {
			if _, ok := m[name]; !ok {
				m[name] = struct{}{}
				unique = append(unique, name)
			}
		}
		g.Package[manager].Names = unique
	}
	return g.save()
}

// AddDotfile add the dotfile config to the Gotfile.
func (g *Gotfile) AddDotfile(src, dest string) error {
	dotfile := &Dotfile{
		Src:  src,
		Dest: dest,
	}
	g.Dotfile = append(g.Dotfile, dotfile)
	return g.save()
}
