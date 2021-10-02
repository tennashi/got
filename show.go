package got

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type ShowCommandConfig struct {
	DataDir string
	IsDebug bool
}

type ShowCommand struct {
	out io.Writer

	repository *InstalledPackageRepository
	prompter   *Prompter
	printer    *TablePrinter
}

func NewShowCommand(ioStream *IOStream, cfg *ShowCommandConfig) (*ShowCommand, error) {
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

	return &ShowCommand{
		out:        ioStream.Out,
		repository: repo,
		prompter:   NewPrompter(ioStream, prompterCfg),
		printer:    printer,
	}, nil
}

func (c *ShowCommand) Run(pkgName string) error {
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

	fmt.Fprintln(c.out, "Package path:", found.Path)
	fmt.Fprintln(c.out, "Version:", found.Version)
	fmt.Fprintln(c.out, "Executables:")
	c.printer.PrintExecutables(found.Executables)

	return nil
}
