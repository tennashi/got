package got

import "path/filepath"

type ListCommandConfig struct {
	DataDir string
	IsDebug bool
}

type ListCommand struct {
	repository *InstalledPackageRepository
	printer    *TablePrinter
}

func NewListCommand(ioStream *IOStream, cfg *ListCommandConfig) (*ListCommand, error) {
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

	return &ListCommand{
		repository: repo,
		printer:    printer,
	}, nil
}

func (c *ListCommand) Run() error {
	pkgs, err := c.repository.List()
	if err != nil {
		return err
	}

	return c.printer.PrintInstalledPackages(pkgs)
}
