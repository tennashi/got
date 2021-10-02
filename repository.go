package got

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"

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

type sortablePackages []InstalledPackage

func (s sortablePackages) Len() int           { return len(s) }
func (s sortablePackages) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortablePackages) Less(i, j int) bool { return s[i].Path < s[j].Path }

func (r *InstalledPackageRepository) List() ([]InstalledPackage, error) {
	r.debugL.Printf("start (*InstalledPackageRepository).List()\n")

	allEntries, err := r.store.GetAll()
	if err != nil {
		return nil, err
	}

	pkgs := make([]InstalledPackage, 0, len(allEntries))
	for _, e := range allEntries {
		pkg := InstalledPackage{}
		err := json.Unmarshal(e, &pkg)
		if err != nil {
			return nil, err
		}

		pkgs = append(pkgs, pkg)
	}

	sort.Sort(sortablePackages(pkgs))

	r.debugL.Printf("end (*InstalledPackageRepository).List()\n")

	return pkgs, nil
}

func (r *InstalledPackageRepository) Save(pkg *InstalledPackage) error {
	r.debugL.Printf("start (*InstalledPackageRepository).Save(%v)\n", pkg)

	r.debugL.Printf("end (*InstalledPackageRepository).Save(%v)\n", pkg)

	return r.store.Save(string(pkg.Path), pkg)
}

func (r *InstalledPackageRepository) Delete(pkg *InstalledPackage) error {
	r.debugL.Printf("start (*InstalledPackageRepository).Save(%v)\n", pkg)

	r.debugL.Printf("end (*InstalledPackageRepository).Save(%v)\n", pkg)

	return r.store.Delete(string(pkg.Path))
}
