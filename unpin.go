package got

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type UnpinCommandConfig struct {
	DataDir string

	IsDebug bool
}

type UnpinCommand struct {
	out io.Writer

	repository *InstalledPackageRepository
	prompter   *Prompter
	printer    *TablePrinter
}

func NewUnpinCommand(ioStream *IOStream, cfg *UnpinCommandConfig) (*UnpinCommand, error) {
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

	return &UnpinCommand{
		out:        ioStream.Out,
		repository: repo,
		printer:    printer,
		prompter:   NewPrompter(ioStream, prompterCfg),
	}, nil
}

func (c *UnpinCommand) Run(pkgName string) error {
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

	found.IsPinned = false

	if err := c.repository.Save(found); err != nil {
		return err
	}

	c.printer.PrintInstalledPackages([]InstalledPackage{*found})

	return nil
}
