package got_test

import (
	"path/filepath"
	"reflect"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	got "github.com/tennashi/got/lib"
)

func TestFilepath_ExpandPath(t *testing.T) {
	home, _ := homedir.Dir()

	cases := map[string]struct {
		input string
		want  string
		err   bool
	}{
		"valid": {
			input: "/etc",
			want:  "/etc",
			err:   false,
		},
		"homedir expand": {
			input: "~/hoge",
			want:  filepath.Join(home, "hoge"),
			err:   false,
		},
		"use tilda outside the beginning": {
			input: "~/hoge/~",
			want:  filepath.Join(home, "hoge", "~"),
			err:   false,
		},
		"use ..": {
			input: "~/..",
			want:  "/home",
			err:   false,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			get, err := got.ExpandPath(tt.input)
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
