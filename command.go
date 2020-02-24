package got

import (
	"os"
	"path/filepath"
)

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
