package got

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tennashi/got/io"
	app_path "github.com/tennashi/got/path"
)

func get(ioStream *io.Stream, pkgName, cmdName string, isUpdate bool) error {
	pkgName = fullPackageName(pkgName, cmdName)
	dataDir, err := app_path.EnsureDataDir()
	if err != nil {
		return err
	}
	cmdCtx := newCommandContext(ioStream, dataDir)

	if err := ensureGoModInit(cmdCtx, "got_local"); err != nil {
		return err
	}
	if err := goGet(cmdCtx, pkgName, isUpdate); err != nil {
		return err
	}
	if err := appendImport(dataDir, pkgName); err != nil {
		return err
	}
	if err := goModTidy(cmdCtx); err != nil {
		return err
	}
	return nil
}

func getAll(ioStream *io.Stream, isUpdate bool) error {
	dataDir, err := app_path.EnsureDataDir()
	if err != nil {
		return err
	}

	imports, err := getImports(dataDir)
	if err != nil {
		return err
	}

	for _, path := range imports {
		if err := get(ioStream, path, "", isUpdate); err != nil {
			return err
		}
	}
	return nil
}

func goGet(cmdCtx *commandContext, pkgName string, isUpdate bool) error {
	args := []string{"get"}
	if isUpdate {
		args = append(args, "-u")
	}
	args = append(args, pkgName)
	return cmdCtx.Exec("go", args)
}

func ensureGoModInit(cmdCtx *commandContext, modName string) error {
	if _, err := os.Stat(filepath.Join(cmdCtx.workDir, "go.mod")); err != nil {
		args := []string{"mod", "init", modName}
		return cmdCtx.Exec("go", args)

	}
	return nil
}

func goModTidy(cmdCtx *commandContext) error {
	args := []string{"mod", "tidy"}
	return cmdCtx.Exec("go", args)
}

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
