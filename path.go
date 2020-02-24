package got

import (
	"os"
	"path/filepath"
	"runtime"
)

func ensureDataDir() (string, error) {
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(os.Getenv("APPDATA"), "got")
	} else if runtime.GOOS == "darwin" {
		path = filepath.Join("/Library", "Application Support", "got")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, ".local", "share", "got")
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	return path, nil
}

func goBinDir() (string, error) {
	path := os.Getenv("GOBIN")
	if path != "" {
		return path, nil
	}
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		return filepath.Join(goPath, "bin"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "go", "bin"), nil
}
