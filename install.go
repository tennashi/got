package got

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
)

type InstallCommandConfig struct {
	DataDir           string
	BinDir            string
	InstallAllCommand bool

	IsDebug bool
}

type InstallCommand struct {
	installAllCommand bool

	out io.Writer

	installer  *PackageInstaller
	repository *InstalledPackageRepository
	linker     *ExecutableLinker
	prompter   *Prompter
}

func NewInstallCommand(ioStream *IOStream, cfg *InstallCommandConfig) (*InstallCommand, error) {
	installerCfg := &PackageInstallerConfig{
		BaseDir: cfg.DataDir,
		IsDebug: cfg.IsDebug,
	}
	installer, err := NewPackageInstaller(ioStream, installerCfg)
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

	prompterCfg := &PrompterConfig{
		IsDebug: cfg.IsDebug,
	}

	return &InstallCommand{
		installAllCommand: cfg.InstallAllCommand,
		out:               ioStream.Out,
		installer:         installer,
		repository:        repo,
		linker:            linker,
		prompter:          NewPrompter(ioStream, prompterCfg),
	}, nil
}

func (c *InstallCommand) Run(pkgName string) error {
	installPkg, err := NewInstallPackage(pkgName, c.installAllCommand)
	if err != nil {
		return err
	}

	nfErr := &PackageNotFoundError{}
	if _, err = c.repository.Get(installPkg.Path); !errors.As(err, &nfErr) {
		if err == nil {
			fmt.Fprintf(c.out, "This packages already installed: %s\n", installPkg.Path)

			return &AlreadyInstalledError{Path: installPkg.Path}
		}

		return err
	}

	fmt.Fprintf(c.out, "This packages will be installed: %s\n", installPkg.Path)

	installedPkg, err := c.installer.Install(installPkg)
	if err != nil {
		return err
	}

	c.prompter.SelectExecutableToDisable(installedPkg)

	for _, exec := range installedPkg.Executables {
		fmt.Fprintf(c.out, "Installed the executable: %s\n", exec.Path)
		fmt.Fprintf(c.out, "Linking the executable: %s\n", exec.Name)

		if err := c.linker.Link(exec); err != nil {
			aeErr := &AlreadyExistsError{}
			if !errors.As(err, &aeErr) {
				return err
			}

			isOverwrite := c.prompter.ChooseToForceOverwrite(aeErr.Path)

			if isOverwrite {
				fmt.Fprintf(c.out, "Force linking the executable: %s\n", exec.Name)
				err := c.linker.ForceLink(exec)
				if err != nil {
					return err
				}
			}
		}
	}

	if err := c.repository.Save(installedPkg); err != nil {
		return err
	}

	return nil
}
