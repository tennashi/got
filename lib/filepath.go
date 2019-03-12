package lib

import (
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func ExpandPath(path string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	path = strings.Replace(path, "~", home, 1)
	return filepath.Abs(path)
}
