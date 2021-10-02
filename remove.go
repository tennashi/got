package got

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type RemoveCommandConfig struct {
	DataDir string
	BinDir  string

	IsDebug bool
}

type RemoveCommand struct {
	out io.Writer

	repository  *InstalledPackageRepository
	uninstaller *PackageUninstaller
	linker      *ExecutableLinker
	prompter    *Prompter
	printer     *TablePrinter
}

func NewRemoveCommand(ioStream *IOStream, cfg *RemoveCommandConfig) (*RemoveCommand, error) {
	uninstallerCfg := &PackageUninstallerConfig{
		BaseDir: cfg.DataDir,
		IsDebug: cfg.IsDebug,
	}
	uninstaller, err := NewPackageUninstaller(ioStream, uninstallerCfg)
	if err != nil {
		return nil, err
	}

	linkerCfg := &ExecutableLinkerConfig{
		BinDir:  cfg.BinDir,
		IsDebug: cfg.IsDebug,
	}
	linker, err := NewExecutableLinker(ioStream, linkerCfg)
	if err != nil {
		return nil, err
	}

	repoCfg := &InstalledPackageRepositoryConfig{
		FilePath: filepath.Join(cfg.DataDir, "package.lock.json"),
		IsDebug:  cfg.IsDebug,
	}
	repo, err := NewInstalledPackageRepository(ioStream, repoCfg)
	if err != nil {
		return nil, err
	}

	printerCfg := &TablePrinterConfig{
		IsDebug: cfg.IsDebug,
	}
	printer := NewTablePrinter(ioStream, printerCfg)

	prompterCfg := &PrompterConfig{
		IsDebug: cfg.IsDebug,
	}

	return &RemoveCommand{
		out:         ioStream.Out,
		repository:  repo,
		uninstaller: uninstaller,
		linker:      linker,
		printer:     printer,
		prompter:    NewPrompter(ioStream, prompterCfg),
	}, nil
}

func (c *RemoveCommand) Run(pkgName string) error {
	pkgs, err := c.repository.List()
	if err != nil {
		return err
	}

	candidate := make([]InstalledPackage, 0)
	for _, pkg := range pkgs {
		if strings.HasSuffix(string(pkg.Path), pkgName) {
			candidate = append(candidate, pkg)
		}
	}

	if len(candidate) == 0 {
		return &PackageNotFoundError{Path: PackagePath(pkgName)}
	}

	found := &candidate[0]
	if len(candidate) > 1 {
		fmt.Fprintln(c.out, "Multiple installed packages with the specified name were found.")
		found = c.prompter.SelectPackage(candidate)
	}

	for _, exec := range found.Executables {
		exec.Disable = true
		if err := c.linker.Unlink(exec); err != nil {
			return err
		}
	}

	if err := c.uninstaller.Uninstall(found); err != nil {
		return err
	}

	if err := c.repository.Delete(found); err != nil {
		return err
	}

	return nil
}
