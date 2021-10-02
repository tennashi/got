package got

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/module"
)

type InstalledPackages []InstalledPackage

func (p InstalledPackages) UpgradeTargets(isAll bool) []InstallPackage {
	installPkgs := []InstallPackage{}
	for _, pkg := range p {
		if !pkg.IsPinned {
			installPkgs = append(installPkgs, InstallPackage{
				Path:    pkg.Path,
				Version: "latest",
				IsAll:   isAll,
			})
		}
	}

	return installPkgs
}

type InstalledPackage struct {
	Path        PackagePath
	Version     string
	IsPinned    bool
	Executables []*Executable
}

type Executable struct {
	Name    string
	Path    string
	Disable bool
}

type ParsePackageError struct {
	Target string
	Err    error
}

func (e *ParsePackageError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Err == nil {
		return fmt.Sprintf("failed to parse package: %s", e.Target)
	}

	return fmt.Sprintf("failed to parse package: %s: %s", e.Target, e.Err.Error())
}

func (e *ParsePackageError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

type InstallPackages []InstallPackage

func (p InstallPackages) Pathes() []PackagePath {
	pathes := make([]PackagePath, 0, len(p))
	for _, pkg := range p {
		pathes = append(pathes, pkg.Path)
	}

	return pathes
}

type InstallPackage struct {
	Path    PackagePath
	Version string
	IsAll   bool
}

func NewInstallPackage(rawPkgName string, isAll bool) (*InstallPackage, error) {
	pair := strings.SplitN(rawPkgName, "@", 2)

	if len(pair) == 0 {
		return nil, &ParsePackageError{
			Target: rawPkgName,
			Err:    errors.New("split into 0 pieces"),
		}
	}

	pkgPath, err := NewPackagePath(pair[0])
	if err != nil {
		return nil, &ParsePackageError{
			Target: rawPkgName,
			Err:    err,
		}
	}

	if len(pair) == 1 {
		return &InstallPackage{
			Path:    pkgPath,
			Version: "latest",
			IsAll:   isAll,
		}, nil
	}

	return &InstallPackage{
		Path:    pkgPath,
		Version: pair[1],
		IsAll:   isAll,
	}, nil
}

func (p *InstallPackage) String() string {
	if p == nil {
		return ""
	}

	if p.IsAll {
		return string(p.Path) + "/...@" + p.Version
	}

	return string(p.Path) + "@" + p.Version
}

type PackagePath string

func NewPackagePath(rawPkgPath string) (PackagePath, error) {
	err := module.CheckPath(rawPkgPath)
	if err != nil {
		if !strings.Contains(err.Error(), "missing dot in first path element") {
			return "", &ParsePackageError{
				Target: rawPkgPath,
				Err:    err,
			}
		}

		return PackagePath("github.com/" + rawPkgPath), nil
	}

	return PackagePath(rawPkgPath), nil
}
