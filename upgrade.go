package got

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type UpgradeCommandConfig struct {
	InstallAllCommand bool
	DataDir           string
	BinDir            string

	IsDebug bool
}

type UpgradeCommand struct {
	installAllCommand bool

	out io.Writer

	installer  *PackageInstaller
	repository *InstalledPackageRepository
	linker     *ExecutableLinker
	prompter   *Prompter
}

func NewUpgradeCommand(ioStream *IOStream, cfg *UpgradeCommandConfig) (*UpgradeCommand, error) {
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

	return &UpgradeCommand{
		installAllCommand: cfg.InstallAllCommand,
		out:               ioStream.Out,
		installer:         installer,
		linker:            linker,
		repository:        repo,
		prompter:          NewPrompter(ioStream, prompterCfg),
	}, nil
}

func (c *UpgradeCommand) Run(pkgName string) error {
	allPackages, err := c.repository.List()
	if err != nil {
		return err
	}

	var candidate []InstalledPackage
	if pkgName != "" {

		for _, pkg := range allPackages {
			if strings.HasSuffix(string(pkg.Path), pkgName) {
				candidate = append(candidate, pkg)
			}
		}

		if len(candidate) == 0 {
			return &PackageNotFoundError{Path: PackagePath(pkgName)}
		}

		if len(candidate) > 1 {
			fmt.Fprintln(c.out, "Multiple installed packages with the specified name were found.")
			candidate = []InstalledPackage{*c.prompter.SelectPackage(candidate)}
		}
	} else {
		candidate = allPackages
	}

	upgradeTargets := InstalledPackages(candidate).UpgradeTargets(c.installAllCommand)

	fmt.Fprintln(c.out, "This packages will be upgraded:")
	for _, path := range InstallPackages(upgradeTargets).Pathes() {
		fmt.Fprintf(c.out, "\t%s\n", path)
	}
	for _, t := range upgradeTargets {
		upgradedPkg, err := c.installer.Install(&t)
		if err != nil {
			return err
		}

		currentPkg, err := c.repository.Get(upgradedPkg.Path)

		executables := make([]*Executable, 0, len(upgradedPkg.Executables))
		for _, exec := range upgradedPkg.Executables {
			isNew := true
			for _, currentExec := range currentPkg.Executables {
				if exec.Path == currentExec.Path {
					exec.Disable = currentExec.Disable
					isNew = false
				}
			}

			if isNew {
				fmt.Fprintln(c.out, "Found a new executable")
				exec = c.prompter.SelectExecutableToDisable(exec)
			}

			executables = append(executables, exec)
		}

		upgradedPkg.Executables = executables

		for _, exec := range upgradedPkg.Executables {
			if exec.Disable {
				continue
			}

			fmt.Fprintf(c.out, "Linking the executable: %s\n", exec.Name)

			if err := c.linker.ForceLink(exec); err != nil {
				return err
			}
		}

		if err := c.repository.Save(upgradedPkg); err != nil {
			return err
		}
	}

	return nil
}
