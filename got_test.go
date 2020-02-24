package got

import (
	"errors"
	"io/ioutil"
	"reflect"
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
			testShowHelp(t)
		})
	}
}

func TestGot_showHelp(t *testing.T) {
	cases := map[string]struct {
		input error
	}{
		"no error": {
			input: nil,
		},
		"error occured": {
			input: errors.New("error"),
		},
	}

	r := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			r.err = tt.input
			got := r.showHelp()
			if !reflect.DeepEqual(got, tt.input) {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.input, got)
			}
			testShowHelp(t)
		})
	}
}

func testShowHelp(t *testing.T) {
	t.Helper()
	out, _ := ioutil.ReadAll(&testOut)
	if string(out) != help {
		t.Fatal("should show help")
	}
}
