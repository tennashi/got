package got

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tennashi/got/path"
	app_path "github.com/tennashi/got/path"
)

func remove(ioStream *ioStream, targetName string) error {
	cmdName := filepath.Base(targetName)

	dataDir, err := app_path.EnsureDataDir()
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
	goBinDir, err := path.GoBinDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(goBinDir, cmdName))
}
