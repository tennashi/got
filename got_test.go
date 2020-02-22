package got

import (
	"io/ioutil"
	"testing"
)

func TestGot_parse_S(t *testing.T) {
	cases := map[string]struct {
		input []string
	}{
		"sub command only": {
			input: []string{"version"},
		},
		"sub command args": {
			input: []string{"version", "hoge", "fuga"},
		},
	}

	r := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			err := r.parse(tt.input)
			if err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if _, ok := r.curSubCmd.(*versionCmd); !ok {
				t.Fatal("should be version command")
			}
		})
	}
}

func TestGot_parse_F(t *testing.T) {
	cases := map[string]struct {
		input []string
		err   bool
	}{
		"empty": {
			input: []string{},
		},
		"nil": {
			input: nil,
		},
		"invalid command name": {
			input: []string{"invalid"},
		},
		"invalid flag": {
			input: []string{"-invalid", "version"},
		},
	}

	r := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			err := r.parse(tt.input)
			if err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			checkShowHelp(t)
		})
	}
}

func checkShowHelp(t *testing.T) {
	t.Helper()
	out, _ := ioutil.ReadAll(&testOut)
	if string(out) != help {
		t.Fatal("should show help")
	}
}
