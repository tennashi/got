package got

import (
	"errors"
	"fmt"
	"log"

	"github.com/schollz/jsonstore"
)

type PackageNotFoundError struct {
	Path PackagePath
}

func (e *PackageNotFoundError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("the package is not found: %s", e.Path)
}

type InstalledPackageRepositoryConfig struct {
	FilePath string
	IsDebug  bool
}

type InstalledPackageRepository struct {
	store *JSONStore

	debugL *log.Logger
}

func NewInstalledPackageRepository(ioStream *IOStream, cfg *InstalledPackageRepositoryConfig) (*InstalledPackageRepository, error) {
	storeCfg := &JSONStoreConfig{
		FilePath: cfg.FilePath,
		IsDebug:  cfg.IsDebug,
	}
	store, err := NewJSONStore(ioStream, storeCfg)
	if err != nil {
		return nil, err
	}

	return &InstalledPackageRepository{
		store:  store,
		debugL: NewDebugLogger(ioStream.Err, "repository", cfg.IsDebug),
	}, nil
}

func (r *InstalledPackageRepository) Get(pkgPath PackagePath) (*InstalledPackage, error) {
	r.debugL.Printf("start (*InstalledPackageRepository).Get(%s)\n", pkgPath)

	installedPkg := &InstalledPackage{}
	if err := r.store.Get(string(pkgPath), installedPkg); err != nil {
		var nskErr jsonstore.NoSuchKeyError
		if errors.As(err, &nskErr) {
			return nil, &PackageNotFoundError{Path: pkgPath}
		}

		return nil, err
	}

	r.debugL.Printf("end (*InstalledPackageRepository).Get(%s)\n", pkgPath)

	return installedPkg, nil
}

func (r *InstalledPackageRepository) Save(pkg *InstalledPackage) error {
	r.debugL.Printf("start (*InstalledPackageRepository).Save(%v)\n", pkg)

	r.debugL.Printf("end (*InstalledPackageRepository).Save(%v)\n", pkg)

	return r.store.Save(string(pkg.Path), pkg)
}
