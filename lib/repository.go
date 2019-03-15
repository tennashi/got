package lib

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Repository interface {
	Clone()
}

type Git struct {
	URL    string
	Target string
}

func NewGit(url, target string) *Git {
	return &Git{
		URL:    url,
		Target: target,
	}
}

func (g *Git) Clone() error {
	if _, err := os.Stat(g.Target); err == nil {
		return errors.New("exist dotfile repository")
	}
	out, err := exec.Command("git", "clone", g.URL, g.Target).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}
