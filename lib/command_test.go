package got_test

import (
	"bytes"
	"testing"

	got "github.com/tennashi/got/lib"
)

func TestCommand_Run(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  string
		err   bool
	}{
		"valid": {
			input: []string{"echo", "hoge"},
			want:  "hoge\n",
			err:   false,
		},
	}

	for caseName, tt := range cases {
		stdin := bytes.NewBufferString("")
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		testCommand := &got.Command{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		}
		t.Run(caseName, func(t *testing.T) {
			err := testCommand.Run(tt.input[0], tt.input[1:]...)
			get := stdout.String()
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if get != tt.want {
				t.Fatalf("\n\tgot: %v\n\twant: %v", get, tt.want)
			}
		})
	}
}
