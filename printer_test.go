package got_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestTablePrinter_PrintInstallPackages(t *testing.T) {
	cases := []struct {
		input []got.InstalledPackage
		want  []string
		err   bool
	}{
		{
			input: []got.InstalledPackage{
				{
					Path:    "github.com/tennashi/got",
					Version: "latest",
					Executables: []*got.Executable{
						{
							Name:    "got",
							Path:    "/path/to/got",
							Disable: false,
						},
					},
				},
				{
					Path:    "github.com/tennashi/got-2",
					Version: "latest",
					Executables: []*got.Executable{
						{
							Name:    "got-2-1",
							Path:    "/path/to/got-2-1",
							Disable: false,
						},
						{
							Name:    "got-2-2",
							Path:    "/path/to/got-2-2",
							Disable: false,
						},
					},
				},
				{
					Path:    "github.com/tennashi/got-3",
					Version: "",
					Executables: []*got.Executable{
						{
							Name:    "got-3",
							Path:    "/path/to/got-3",
							Disable: false,
						},
					},
				},
			},
			want: []string{
				"NAME                       VERSION  EXECUTABLES  ",
				"github.com/tennashi/got    latest   got",
				"github.com/tennashi/got-2  latest   got-2-1,got-2-2",
				"github.com/tennashi/got-3           got-3",
				"",
			},
			err: false,
		},
		{
			input: []got.InstalledPackage{},
			want: []string{
				"NAME  VERSION  EXECUTABLES  ",
				"",
			},
			err: false,
		},
		{
			input: nil,
			want: []string{
				"NAME  VERSION  EXECUTABLES  ",
				"",
			},
			err: false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ioStream := newTestIOStream()
			cfg := &got.TablePrinterConfig{
				IsDebug: true,
			}

			p := got.NewTablePrinter(ioStream, cfg)

			err := p.PrintInstalledPackages(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error but not")
			}

			got := ioStream.Out.(*bytes.Buffer).String()
			if diff := cmp.Diff(strings.Join(tt.want, "\n"), got); diff != "" {
				t.Fatalf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
