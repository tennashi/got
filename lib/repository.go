package got

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Repository defines the operation of VCS repository.
type Repository interface {
	Clone()
}

// Git represent the git command.
type Git struct {
	URL    string
	Target string
}

// NewGit initialize Git type from the remote repository URL and the local repository path.
func NewGit(url, target string) *Git {
	return &Git{
		URL:    url,
		Target: target,
	}
}

// Clone clone the git repository from URL to Target.
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
