package got

import (
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// ExpandPath is a wrapper for `filepath.Abs()`.
// This function expand `~` to your home directory like `/home/you`.
func ExpandPath(path string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	path = strings.Replace(path, "~", home, 1)
	return filepath.Abs(path)
}
