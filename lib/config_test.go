package got_test

import (
	"path/filepath"
	"reflect"
	"testing"

	got "github.com/tennashi/got/lib"
)

func TestConfig_InitConfig(t *testing.T) {
	cases := map[string]struct {
		input string
		want  *got.Config
		err   bool
	}{
		"valid": {
			input: filepath.Join("..", "testdata", "config_valid.golden"),
			want:  testConfig["valid"],
			err:   false,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			get, err := got.InitConfig(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(get, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", get, tt.want)
			}
		})
	}
}

func TestConfig_setPaths(t *testing.T) {
	cases := map[string]struct {
		input string
		err   bool
	}{}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			err := got.CsetPaths(testConfig["empty"], tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
		})
	}
}

func TestConfig_addPath(t *testing.T) {
	cases := map[string]struct {
		input string
		want  []string
	}{
		"valid": {
			input: filepath.Join("..", "testdata", "config_valid.golden"),
			want:  []string{filepath.Join("..", "testdata", "config_valid.golden")},
		},
		"unexist path": {
			input: "hogehoge",
			want:  nil,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			got.CaddPath(testConfig["empty"], tt.input)
			get := testConfig["empty"].ExportCpaths()
			if !reflect.DeepEqual(get, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", get, tt.want)
			}
			resetTestConfig()
		})
	}
}
