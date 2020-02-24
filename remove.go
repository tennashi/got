package got

import (
	"os"
	"path/filepath"
	"strings"
)

func remove(ioStream *ioStream, targetName string) error {
	cmdName := filepath.Base(targetName)

	dataDir, err := ensureDataDir()
	if err != nil {
		return err
	}

	var pkgName string
	if !strings.ContainsRune(targetName, '/') {
		pkgName, err = getPackageName(dataDir, targetName)
		if err != nil {
			return err
		}
	} else {
		pkgName = fullPackageName(targetName, "")
	}

	cmdCtx := newCommandContext(ioStream, dataDir)
	if err := removeImport(dataDir, pkgName); err != nil {
		return err
	}
	if err := removeCommand(cmdName); err != nil {
		return err
	}
	return goModTidy(cmdCtx)
}

func removeCommand(cmdName string) error {
	goBinDir, err := goBinDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(goBinDir, cmdName))
}
