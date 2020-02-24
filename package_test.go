package got

import "testing"

func Test_fullPackageName(t *testing.T) {
	type input struct {
		pkgName string
		cmdName string
	}

	cases := map[string]struct {
		input input
		want  string
	}{
		"set full package name and version": {
			input: input{
				pkgName: "github.com/tennashi/got@0.0.1",
				cmdName: "",
			},
			want: "github.com/tennashi/got@0.0.1",
		},
		"set full package name": {
			input: input{
				pkgName: "github.com/tennashi/got",
				cmdName: "",
			},
			want: "github.com/tennashi/got",
		},
		"set short package name and version": {
			input: input{
				pkgName: "tennashi/got@0.0.1",
				cmdName: "",
			},
			want: "github.com/tennashi/got@0.0.1",
		},
		"set short package name": {
			input: input{
				pkgName: "tennashi/got",
				cmdName: "",
			},
			want: "github.com/tennashi/got",
		},
		"set command name with full package name and version": {
			input: input{
				pkgName: "github.com/tennashi/got@0.0.1",
				cmdName: "command",
			},
			want: "github.com/tennashi/got/cmd/command@0.0.1",
		},
		"set command name with full package name": {
			input: input{
				pkgName: "github.com/tennashi/got",
				cmdName: "command",
			},
			want: "github.com/tennashi/got/cmd/command",
		},
		"set command name with short package name and version": {
			input: input{
				pkgName: "tennashi/got@0.0.1",
				cmdName: "command",
			},
			want: "github.com/tennashi/got/cmd/command@0.0.1",
		},
		"set command name with short package name": {
			input: input{
				pkgName: "tennashi/got",
				cmdName: "command",
			},
			want: "github.com/tennashi/got/cmd/command",
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			got := fullPackageName(tt.input.pkgName, tt.input.cmdName)
			if got != tt.want {
				t.Fatalf("%v: want %v, but got: %v", caseName, tt.want, got)
			}
		})
	}
}
