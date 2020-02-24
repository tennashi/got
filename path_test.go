package got

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_goBinDir(t *testing.T) {
	input := "/test/bin"
	caseName := "set $GOBIN, unset $GOPATH"
	t.Run(caseName, func(t *testing.T) {
		want := input

		os.Setenv("GOBIN", input)
		got, err := goBinDir()
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		if got != want {
			t.Fatalf("%v: want %v, but got: %v", caseName, want, got)
		}
	})

	caseName = "unset $GOBIN, set $GOPATH"
	t.Run(caseName, func(t *testing.T) {
		want := filepath.Join(input, "bin")

		os.Unsetenv("GOBIN")
		os.Setenv("GOPATH", input)
		got, err := goBinDir()
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		if got != want {
			t.Fatalf("%v: want %v, but got: %v", caseName, want, got)
		}
	})

	caseName = "unset $GOBIN, unset $GOPATH"
	t.Run(caseName, func(t *testing.T) {
		want := filepath.Join(mustHomeDir(t), "go", "bin")

		os.Unsetenv("GOPATH")
		got, err := goBinDir()
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		if got != want {
			t.Fatalf("%v: want %v, but got: %v", caseName, want, got)
		}
	})

	caseName = "set $GOBIN, set $GOPATH"
	t.Run(caseName, func(t *testing.T) {
		want := input

		os.Setenv("GOBIN", input)
		os.Setenv("GOPATH", input)
		got, err := goBinDir()
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		if got != want {
			t.Fatalf("%v: want %v, but got: %v", caseName, want, got)
		}
	})
}

func mustHomeDir(t *testing.T) string {
	t.Helper()
	path, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path
}
