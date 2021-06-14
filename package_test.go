package got

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInstallPackage(t *testing.T) {
	cases := []struct {
		input string
		want  string
		err   bool
	}{
		{
			input: "github.com/tennashi/got@0.0.1",
			want:  "github.com/tennashi/got@0.0.1",
			err:   false,
		},
		{
			input: "tennashi/got@0.0.1",
			want:  "github.com/tennashi/got@0.0.1",
			err:   false,
		},
		{
			input: "got@0.0.1",
			want:  "github.com/got@0.0.1",
			err:   false,
		},
		{
			input: "github.com/tennashi/got",
			want:  "github.com/tennashi/got@latest",
			err:   false,
		},
		{
			input: "-/tennashi/got",
			want:  "",
			err:   true,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			got, err := NewInstallPackage(tt.input, false)
			if !tt.err && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error but not")
			}
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
