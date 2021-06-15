package got_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestPackageInstaller_Install(t *testing.T) {
	ioStream := newTestIOStream()
	tmpDir := t.TempDir()
	i, err := got.NewPackageInstaller(ioStream, &got.PackageInstallerConfig{
		BaseDir: tmpDir,
		IsDebug: true,
	})
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	input, err := got.NewInstallPackage("tennashi/got", true)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	g, err := i.Install(input)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	want := &got.InstalledPackage{
		Path:    "github.com/tennashi/got",
		Version: "latest",
		Executables: []*got.Executable{
			{
				Name:    "got",
				Path:    filepath.Join(tmpDir, "github.com/tennashi/got/got"),
				Disable: false,
			},
		},
	}

	if diff := cmp.Diff(want, g); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
}
