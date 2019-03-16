package lib

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type SymLink struct {
	Src  string
	Dest string
}

func NewSymLink(dirPath string, dotfile Dotfile) (*SymLink, error) {
	dest, err := ExpandPath(dotfile.Dest)
	if err != nil {
		return nil, err
	}
	symLink := &SymLink{
		Src:  filepath.Join(dirPath, dotfile.Src),
		Dest: dest,
	}
	if err := symLink.check(); err != nil {
		return nil, err
	}
	return symLink, nil
}

func (s *SymLink) check() error {
	if _, err := os.Stat(s.Dest); err == nil {
		return errors.New("exist dest file")
	}

	if info, err := os.Stat(s.Src); err != nil {
		return err
	} else if info.IsDir() {
		s.Src += "/"
	}
	return nil
}

func (s *SymLink) Make() error {
	out, err := exec.Command("ln", "-s", s.Src, s.Dest).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}
