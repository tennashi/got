package got

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type DisableCommandConfig struct {
	DataDir string
	BinDir  string

	IsDebug bool
}

type DisableCommand struct {
	out io.Writer

	repository *InstalledPackageRepository
	linker     *ExecutableLinker
	prompter   *Prompter
	printer    *TablePrinter
}

func NewDisableCommand(ioStream *IOStream, cfg *DisableCommandConfig) (*DisableCommand, error) {
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

	return &DisableCommand{
		out:        ioStream.Out,
		repository: repo,
		linker:     linker,
		printer:    printer,
		prompter:   NewPrompter(ioStream, prompterCfg),
	}, nil
}

func (c *DisableCommand) Run(pkgName, execName string) error {
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

	var targetExec *Executable
	for _, exec := range found.Executables {
		if exec.Name == execName {
			targetExec = exec
		}
	}

	if targetExec == nil {
		return errors.New("not found")
	}

	targetExec.Disable = true

	if err := c.linker.Unlink(targetExec); err != nil {
		return err
	}

	if err := c.repository.Save(found); err != nil {
		return err
	}

	c.printer.PrintInstalledPackages([]InstalledPackage{*found})

	return nil
}
