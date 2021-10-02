package got_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestExecutableLinker_Link(t *testing.T) {
	tmpDir := t.TempDir()

	input := &got.Executable{
		Name:    "got",
		Path:    filepath.Join("testdata/installed_bin", "got"),
		Disable: false,
	}

	ioStream := newTestIOStream()
	cfg := &got.ExecutableLinkerConfig{
		BinDir:  tmpDir,
		IsDebug: true,
	}
	l, err := got.NewExecutableLinker(ioStream, cfg)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	err = l.Link(input)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	dirEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	if diff := cmp.Diff(1, len(dirEntries)); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	info, err := os.Lstat(filepath.Join(tmpDir, "got"))
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	if diff := cmp.Diff(fs.ModeSymlink, info.Mode().Type()); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
}

func TestExecutableLinker_Unlink(t *testing.T) {
	tmpDir := t.TempDir()

	input := &got.Executable{
		Name:    "got",
		Path:    filepath.Join("testdata/installed_bin", "got"),
		Disable: false,
	}

	ioStream := newTestIOStream()
	cfg := &got.ExecutableLinkerConfig{
		BinDir:  tmpDir,
		IsDebug: true,
	}
	l, err := got.NewExecutableLinker(ioStream, cfg)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	err = l.Link(input)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	dirEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	if diff := cmp.Diff(1, len(dirEntries)); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	input.Disable = true

	err = l.Unlink(input)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	dirEntries, err = os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("should not be error but: %v", err)
	}

	if diff := cmp.Diff(0, len(dirEntries)); diff != "" {
		t.Fatalf("mismatch (-want, +got): %s\n", diff)
	}

	t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
}
