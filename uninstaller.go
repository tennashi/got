package got

import (
	"log"
	"os"
	"path/filepath"
)

type PackageUninstallerConfig struct {
	BaseDir string
	IsDebug bool
}

type PackageUninstaller struct {
	baseDir string

	debugL *log.Logger
}

func NewPackageUninstaller(ioStream *IOStream, cfg *PackageUninstallerConfig) (*PackageUninstaller, error) {
	if cfg.BaseDir == "" {
		return nil, &InvalidParamError{Message: "base directory must not be empty"}
	}

	return &PackageUninstaller{
		baseDir: cfg.BaseDir,
		debugL:  NewDebugLogger(ioStream.Err, "uninstaller", cfg.IsDebug),
	}, nil
}

func (u *PackageUninstaller) Uninstall(pkg *InstalledPackage) error {
	u.debugL.Printf("start (*PackageUninstaller).Uninstall(%v)\n", pkg)

	installedPath := filepath.Join(u.baseDir, string(pkg.Path))
	return os.RemoveAll(installedPath)
}
