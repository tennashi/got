package got

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/otiai10/copy"
)

var testDir = "testdata"
var tmpDir = filepath.Join(testDir, "tmp")

func Test_getImports(t *testing.T) {
	cases := map[string]struct {
		input string
		want  []string
		err   bool
	}{
		"empty import": {
			input: "empty_import",
			want:  []string{},
			err:   false,
		},
		"exist imports": {
			input: "normal",
			want: []string{
				"hoge/fuga/piyo",
				"github.com/tennashi/got",
			},
			err: false,
		},
		"broken file": {
			input: "broken",
			want:  nil,
			err:   true,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			toolsPath := filepath.Join(testDir, "tools", tt.input+".go")
			tmpPath := filepath.Join(tmpDir, "tools.go")
			copy.Copy(toolsPath, tmpPath)
			defer os.RemoveAll(tmpDir)

			got, err := getImports(tmpDir)

			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}

func Test_appendImport(t *testing.T) {
	cases := map[string]struct {
		input string
		want  []string
		err   bool
	}{
		"normal": {
			input: "empty_import",
			want:  []string{"hogehoge"},
			err:   false,
		},
		"exist imports": {
			input: "normal",
			want: []string{
				"hoge/fuga/piyo",
				"hogehoge",
				"github.com/tennashi/got",
			},
			err: false,
		},
		"already imported": {
			input: "hogehoge_imported",
			want: []string{
				"hogehoge",
			},
			err: false,
		},
		"broken file": {
			input: "broken",
			want:  nil,
			err:   true,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			toolsPath := filepath.Join(testDir, "tools", tt.input+".go")
			tmpPath := filepath.Join(tmpDir, "tools.go")
			copy.Copy(toolsPath, tmpPath)
			defer os.RemoveAll(tmpDir)

			err := appendImport(tmpDir, "hogehoge")
			got, _ := getImports(tmpDir)

			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}

func Test_removeImport(t *testing.T) {
	cases := map[string]struct {
		input string
		want  []string
		err   bool
	}{
		"normal": {
			input: "hogehoge_imported",
			want:  []string{},
			err:   false,
		},
		"hogehoge doesn't exist": {
			input: "normal",
			want: []string{
				"hoge/fuga/piyo",
				"github.com/tennashi/got",
			},
			err: false,
		},
		"broken file": {
			input: "broken",
			want:  nil,
			err:   true,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			toolsPath := filepath.Join(testDir, "tools", tt.input+".go")
			tmpPath := filepath.Join(tmpDir, "tools.go")
			copy.Copy(toolsPath, tmpPath)
			defer os.RemoveAll(tmpDir)

			err := removeImport(tmpDir, "hogehoge")
			got, _ := getImports(tmpDir)

			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}

func Test_ensureToolsFile(t *testing.T) {
	caseName := "empty directory"
	t.Run(caseName, func(t *testing.T) {
		os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)

		got, err := ensureToolsFile(tmpDir)
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		toolsPath := filepath.Join(tmpDir, "tools.go")
		if got != toolsPath {
			t.Fatalf("%v: want %v, but got: %v", caseName, toolsPath, got)
		}

		f, err := os.Open(toolsPath)
		defer f.Close()
		if err != nil {
			t.Fatalf("%v: should create %v", caseName, toolsPath)
		}
		data, _ := ioutil.ReadAll(f)
		if string(data) != "package main\n" {
			t.Fatalf("%v: want 'package main', but got: %v", caseName, string(data))
		}
	})

	caseName = "already created tools.go"
	t.Run(caseName, func(t *testing.T) {
		dataDir := filepath.Join(testDir, "normal")
		copy.Copy(dataDir, tmpDir)
		defer os.RemoveAll(tmpDir)

		got, err := ensureToolsFile(tmpDir)
		if err != nil {
			t.Fatalf("should not be error for %v but %v", caseName, err)
		}
		toolsPath := filepath.Join(tmpDir, "tools.go")
		if got != toolsPath {
			t.Fatalf("%v: want %v, but got: %v", caseName, toolsPath, got)
		}

		f, err := os.Open(toolsPath)
		defer f.Close()
		if err != nil {
			t.Fatalf("%v: should create %v", caseName, toolsPath)
		}
		data, _ := ioutil.ReadAll(f)
		if string(data[:12]) != "package main" {
			t.Fatalf("%v: want 'package main', but got: %v", caseName, string(data))
		}
	})
}
