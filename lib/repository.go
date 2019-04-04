package got

import (
	"os"
)

// Repository defines the operation of VCS repository.
type Repository interface {
	CloneOrPull()
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

// CloneOrPull clone the git repository from URL to Target.
// When already exists the Target, this function executes git pull.
func (g *Git) CloneOrPull() error {
	c := NewCommand()
	if _, err := os.Stat(g.Target); err == nil {
		err := c.RunInDir(g.Target, "git", "pull", "origin", "master")
		if err != nil {
			return err
		}
		return nil
	}
	err := c.Run("git", "clone", g.URL, g.Target)
	if err != nil {
		return err
	}
	return nil
}
