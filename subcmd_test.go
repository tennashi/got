package got

import "testing"

func TestVersionCmd_parse(t *testing.T) {
	g := newGot(testIOStream)
	r := newVersionCmd(g)
	err := r.parse([]string{"hoge", "fuga"})
	if err != nil {
		t.Fatalf("should not be error but %v", err)
	}
}

func TestHelpCmd_parse(t *testing.T) {
	cases := map[string]struct {
		input []string
		err   bool
	}{
		"nil": {
			input: nil,
			err:   false,
		},
		"empty array": {
			input: []string{},
			err:   false,
		},
		"some string": {
			input: []string{"hoge", "fuga"},
			err:   false,
		},
		"invalid option": {
			input: []string{"-invalid", "fuga"},
			err:   true,
		},
	}

	g := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			r := newHelpCmd(g)
			err := r.parse(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
		})
	}
}

func TestGetCmd_parse(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  string
		err   bool
	}{
		"some string": {
			input: []string{"hoge", "fuga"},
			want:  "hoge",
			err:   false,
		},
		"list option": {
			input: []string{"-l", "hoge", "fuga"},
			want:  "hoge",
			err:   false,
		},
		"update option": {
			input: []string{"-u", "hoge", "fuga"},
			want:  "hoge",
			err:   false,
		},
		"command name option": {
			input: []string{"-c", "hoge", "fuga"},
			want:  "fuga",
			err:   false,
		},
		"nil": {
			input: nil,
			want:  "",
			err:   true,
		},
		"empty array": {
			input: []string{},
			want:  "",
			err:   true,
		},
		"invalid option": {
			input: []string{"-invalid", "fuga"},
			want:  "",
			err:   true,
		},
	}

	g := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			r := newGetCmd(g)
			err := r.parse(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			got := r.pkgName
			if got != tt.want {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}

func TestGetCmd_parse_cmdNameOption(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  string
		err   bool
	}{
		"provide command name option": {
			input: []string{"-c", "hoge", "fuga"},
			want:  "hoge",
			err:   false,
		},
		"not provide package name": {
			input: []string{"-c", "hoge"},
			want:  "hoge",
			err:   true,
		},
		"provide flag only": {
			input: []string{"-c"},
			want:  "",
			err:   true,
		},
	}

	g := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			r := newGetCmd(g)
			err := r.parse(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			got := r.cmdName
			if got != tt.want {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}

func TestRemoveCmd_parse(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  string
		err   bool
	}{
		"some string": {
			input: []string{"hoge", "fuga"},
			want:  "hoge",
			err:   false,
		},
		"nil": {
			input: nil,
			want:  "",
			err:   true,
		},
		"empty array": {
			input: []string{},
			want:  "",
			err:   true,
		},
		"invalid option": {
			input: []string{"-invalid", "fuga"},
			want:  "",
			err:   true,
		},
	}

	g := newGot(testIOStream)
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			r := newRemoveCmd(g)
			err := r.parse(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			got := r.targetName
			if got != tt.want {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}
