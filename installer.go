package got

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PackageInstallerConfig struct {
	BaseDir string
	IsDebug bool
}

type PackageInstaller struct {
	baseDir  string
	executor *Executor

	debugL *log.Logger
}

func NewPackageInstaller(ioStream *IOStream, cfg *PackageInstallerConfig) (*PackageInstaller, error) {
	if cfg.BaseDir == "" {
		return nil, &InvalidParamError{Message: "base directory must not be empty"}
	}

	executorCfg := &ExecutorConfig{
		IsDebug: cfg.IsDebug,
	}

	return &PackageInstaller{
		baseDir:  cfg.BaseDir,
		executor: NewExecutor(ioStream, executorCfg),
		debugL:   NewDebugLogger(ioStream.Err, "installer", cfg.IsDebug),
	}, nil
}

func (c *PackageInstaller) Install(pkg *InstallPackage) (*InstalledPackage, error) {
	c.debugL.Printf("start (*PackageInstaller).Install(%v)\n", pkg)

	installPath := filepath.Join(c.baseDir, string(pkg.Path))
	c.executor.SetEnv("GOBIN", installPath)
	c.debugL.Printf("destination path: %s\n", installPath)

	args := []string{"install", pkg.String()}

	if err := c.executor.Exec("go", args); err != nil {
		c.debugL.Printf("error occurred in c.executor.Exec(): %v\n", err)
		return nil, err
	}

	dirEntries, err := os.ReadDir(installPath)
	if err != nil {
		c.debugL.Printf("error occurred in os.ReadDir(): %v\n", err)
		return nil, err
	}

	executables := make([]*Executable, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		executables = append(executables, &Executable{
			Name: dirEntry.Name(),
			Path: filepath.Join(installPath, dirEntry.Name()),
		})
		c.debugL.Printf("installed: %s\n", dirEntry.Name())
	}

	version := pkg.Version
	if len(executables) != 0 {
		stdout, err := c.executor.ExecBackground("go", []string{"version", "-m", executables[0].Path})
		if err != nil {
			c.debugL.Printf("error occurred in c.executor.ExecBackground(): %v\n", err)
			return nil, err
		}

		for _, line := range strings.Split(stdout, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "mod") {
				version = strings.Fields(line)[2]
			}
		}
	}

	c.debugL.Printf("end (*PackageInstaller).Install(%v)\n", pkg)
	return &InstalledPackage{
		Path:        pkg.Path,
		Version:     version,
		Executables: executables,
	}, nil
}
