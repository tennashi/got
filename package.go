package got

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
)

func fullPackageName(pkgName, cmdName string) string {
	parts := strings.Split(pkgName, "@")
	repoParts := strings.Split(parts[0], "/")

	var packParts []string
	if !strings.ContainsRune(repoParts[0], '.') {
		packParts = append(packParts, "github.com")
	}
	packParts = append(packParts, repoParts...)
	if cmdName != "" {
		packParts = append(packParts, "cmd", cmdName)
	}
	if len(parts) == 2 {
		packParts[len(packParts)-1] += "@" + parts[1]
	}
	return path.Join(packParts...)
}

func getPackageName(dataDir, cmdName string) (string, error) {
	imports, err := getImports(dataDir)
	if err != nil {
		return "", err
	}
	for _, impPath := range imports {
		if filepath.Base(impPath) == cmdName {
			return impPath, nil
		}
	}
	return "", errors.New("package not found")
}
