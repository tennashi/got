package repository

import (
	"path/filepath"

	"github.com/tennashi/got/path"
)

var targetDir string

func init() {
	dataDir, err := path.EnsureDataDir()
	if err != nil {
		panic(err)
	}
	targetDir = filepath.Join(dataDir, "dotfiles")
}

func (r *Repository) Clone() error {
	args := []string{"clone", r.URI, targetDir}

	if err := r.Cmd(args).Run(); err != nil {
		return err
	}
	return nil
}
