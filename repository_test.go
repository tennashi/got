package got_test

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/schollz/jsonstore"
	"github.com/tennashi/got"
)

var testInstalledPackages = []got.InstalledPackage{
	{
		Path:    "github.com/tennashi/got",
		Version: "latest",
		Executables: []*got.Executable{
			{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: false,
			},
		},
	},
	{
		Path:    "github.com/tennashi/got-2",
		Version: "latest",
		Executables: []*got.Executable{
			{
				Name:    "got-2",
				Path:    "/path/to/got-2",
				Disable: false,
			},
		},
	},
}

func initTestJSONStore(t *testing.T, pkgs []got.InstalledPackage) string {
	t.Helper()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "installed.json")

	store := new(jsonstore.JSONStore)

	for _, pkg := range pkgs {
		err := store.Set(string(pkg.Path), pkg)
		if err != nil {
			t.Fatalf("error occurred in initTestPackageLockFile(): %v", err)
		}
	}

	err := jsonstore.Save(store, path)
	if err != nil {
		t.Fatalf("error occurred in initTestPackageLockFile(): %v", err)
	}

	return path
}

func TestInstalledPackageRepository_Get(t *testing.T) {
	cases := []struct {
		input got.PackagePath
		want  *got.InstalledPackage
		err   bool
	}{
		{
			input: "github.com/tennashi/got",
			want: &got.InstalledPackage{
				Path:    "github.com/tennashi/got",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "got",
						Path:    "/path/to/got",
						Disable: false,
					},
				},
			},
			err: false,
		},
		{
			input: "not-exist",
			want:  nil,
			err:   true,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			jsonPath := initTestJSONStore(t, testInstalledPackages)
			ioStream := newTestIOStream()
			cfg := &got.InstalledPackageRepositoryConfig{
				FilePath: jsonPath,
				IsDebug:  true,
			}
			r, err := got.NewInstalledPackageRepository(ioStream, cfg)
			if err != nil {
				t.Fatalf("should not be error but: %v", err)
			}

			got, err := r.Get(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}

			if tt.err && err == nil {
				t.Fatal("should be error but not")
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
		})
	}
}

func TestInstalledPackageRepository_Save(t *testing.T) {
	cases := []struct {
		input *got.InstalledPackage
		want  *got.InstalledPackage
		err   bool
	}{
		{
			input: &got.InstalledPackage{
				Path:    "github.com/tennashi/got",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "updated",
						Path:    "/path/to/updated",
						Disable: true,
					},
				},
			},
			want: &got.InstalledPackage{
				Path:    "github.com/tennashi/got",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "updated",
						Path:    "/path/to/updated",
						Disable: true,
					},
				},
			},
			err: false,
		},
		{
			input: &got.InstalledPackage{
				Path:    "github.com/tennashi/new",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "new",
						Path:    "/path/to/new",
						Disable: true,
					},
					{
						Name:    "new-2",
						Path:    "/path/to/new-2",
						Disable: false,
					},
				},
			},
			want: &got.InstalledPackage{
				Path:    "github.com/tennashi/new",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "new",
						Path:    "/path/to/new",
						Disable: true,
					},
					{
						Name:    "new-2",
						Path:    "/path/to/new-2",
						Disable: false,
					},
				},
			},
			err: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			jsonPath := initTestJSONStore(t, testInstalledPackages)
			ioStream := newTestIOStream()
			cfg := &got.InstalledPackageRepositoryConfig{
				FilePath: jsonPath,
				IsDebug:  true,
			}
			r, err := got.NewInstalledPackageRepository(ioStream, cfg)
			if err != nil {
				t.Fatalf("should not be error but: %v", err)
			}

			err = r.Save(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}

			if tt.err && err == nil {
				t.Fatal("should be error but not")
			}

			got, err := r.Get(tt.input.Path)
			if err != nil {
				t.Fatalf("should not be error but: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
		})
	}
}

func TestInstalledPackageRepository_List(t *testing.T) {
	jsonPath := initTestJSONStore(t, testInstalledPackages)
	ioStream := newTestIOStream()
	cfg := &got.InstalledPackageRepositoryConfig{
		FilePath: jsonPath,
		IsDebug:  true,
	}
	r, err := got.NewInstalledPackageRepository(ioStream, cfg)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	got, err := r.List()
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	if diff := cmp.Diff(testInstalledPackages, got); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
}

func TestInstalledPackageRepository(t *testing.T) {
	jsonPath := initTestJSONStore(t, nil)

	ioStream := newTestIOStream()
	cfg := &got.InstalledPackageRepositoryConfig{
		FilePath: jsonPath,
		IsDebug:  true,
	}

	r, err := got.NewInstalledPackageRepository(ioStream, cfg)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	t.Run("r.Get() returns a PackageNotFoundError when the store is empty", func(t *testing.T) {
		g, err := r.Get("not exist")
		if err == nil {
			t.Fatal("should be error but not")
		}
		if nfErr := new(got.PackageNotFoundError); !errors.As(err, &nfErr) {
			t.Fatalf("want: *got.PackageNotFoundError, got: %T", err)
		}
		if g != nil {
			t.Fatalf("the return value should be empty but: %v", g)
		}

	})

	t.Run("r.List() returns an empty list", func(t *testing.T) {
		g, err := r.List()
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		if diff := cmp.Diff([]got.InstalledPackage{}, g); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}
	})

	testInstalledPackage := &got.InstalledPackage{
		Path:    "github.com/tennashi/got",
		Version: "latest",
		Executables: []*got.Executable{
			{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: false,
			},
		},
	}

	t.Run("r.Save() will succeed when the store is empty", func(t *testing.T) {
		err := r.Save(testInstalledPackage)
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
	})

	t.Run("r.Get() returns what you have saved", func(t *testing.T) {
		g, err := r.Get(testInstalledPackage.Path)
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		if diff := cmp.Diff(testInstalledPackage, g); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}
	})

	t.Run("r.List() returns saved packages", func(t *testing.T) {
		g, err := r.List()
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		if diff := cmp.Diff([]got.InstalledPackage{*testInstalledPackage}, g); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}
	})

	testInstalledPackage.Version = "updated"

	t.Run("r.Save() will also succeed if path already exists", func(t *testing.T) {
		err := r.Save(testInstalledPackage)
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
	})

	t.Run("r.Get() will return the latest updated one", func(t *testing.T) {
		g, err := r.Get(testInstalledPackage.Path)
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		if diff := cmp.Diff(testInstalledPackage, g); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}
	})

	t.Run("r.List() returns saved packages", func(t *testing.T) {
		g, err := r.List()
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		if diff := cmp.Diff([]got.InstalledPackage{*testInstalledPackage}, g); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}
	})

	t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
}
