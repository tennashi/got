package got

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// SymLink represent symbolic link.
type SymLink struct {
	Src  string
	Dest string
}

// NewSymLink initialize SymLink.
// The dest is absolute path of dotfile.Dest.
// The src is /path/to/your/dotfiles/dotfile.Src.
// This function returns an error if the file exists in the dest path
// or the file does not exist in the src path.
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

// Make symbolic links the src to the dest.
func (s *SymLink) Make() error {
	c := NewCommand()
	err := c.SURun("ln", "-s", s.Src, s.Dest)
	if err != nil {
		return err
	}
	return nil
}
