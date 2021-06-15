package got_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestPrompter_SelectExecutableToDisable(t *testing.T) {
	cases := []struct {
		input *got.InstalledPackage
		stdin string
		want  *got.InstalledPackage
	}{
		{
			input: &got.InstalledPackage{
				Path:    "github.com/tennashi/got",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "got",
						Path:    "/path/to/got",
						Disable: false,
					},
					{
						Name:    "got-2",
						Path:    "/path/to/got-2",
						Disable: false,
					},
				},
			},
			stdin: "y\nn",
			want: &got.InstalledPackage{
				Path:    "github.com/tennashi/got",
				Version: "latest",
				Executables: []*got.Executable{
					{
						Name:    "got",
						Path:    "/path/to/got",
						Disable: false,
					},
					{
						Name:    "got-2",
						Path:    "/path/to/got-2",
						Disable: true,
					},
				},
			},
		},
		{
			input: &got.InstalledPackage{
				Path:        "github.com/tennashi/got",
				Version:     "latest",
				Executables: []*got.Executable{},
			},
			stdin: "",
			want: &got.InstalledPackage{
				Path:        "github.com/tennashi/got",
				Version:     "latest",
				Executables: []*got.Executable{},
			},
		},
		{
			input: &got.InstalledPackage{
				Path:        "github.com/tennashi/got",
				Version:     "latest",
				Executables: nil,
			},
			stdin: "",
			want: &got.InstalledPackage{
				Path:        "github.com/tennashi/got",
				Version:     "latest",
				Executables: nil,
			},
		},
		{
			input: nil,
			stdin: "",
			want:  nil,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ioStream := newTestIOStream()
			ioStream.In.(*bytes.Buffer).WriteString(tt.stdin)

			cfg := &got.PrompterConfig{
				IsDebug: true,
			}
			p := got.NewPrompter(ioStream, cfg)
			p.SelectExecutableToDisable(tt.input)

			if diff := cmp.Diff(tt.want, tt.input); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
		})
	}
}

func TestPrompter_ChooseToForceOverwrite(t *testing.T) {
	cases := []struct {
		input string
		stdin string
		want  bool
	}{
		{
			input: "/path/to/overwrite",
			stdin: "y",
			want:  true,
		},
		{
			input: "/path/to/overwrite",
			stdin: "n",
			want:  false,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ioStream := newTestIOStream()
			ioStream.In.(*bytes.Buffer).WriteString(tt.stdin)

			cfg := &got.PrompterConfig{
				IsDebug: true,
			}
			p := got.NewPrompter(ioStream, cfg)
			got := p.ChooseToForceOverwrite(tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
		})
	}
}

func TestPrompter_AskYN(t *testing.T) {
	cases := []struct {
		input   string
		msg     string
		wantMsg string
		want    bool
	}{
		{
			input:   "y",
			msg:     "test",
			wantMsg: "test [Y/n]: ",
			want:    true,
		},
		{
			input:   "Y",
			msg:     "test",
			wantMsg: "test [Y/n]: ",
			want:    true,
		},
		{
			input:   "n",
			msg:     "test",
			wantMsg: "test [Y/n]: ",
			want:    false,
		},
		{
			input:   "N",
			msg:     "test",
			wantMsg: "test [Y/n]: ",
			want:    false,
		},
		{
			input:   "",
			msg:     "test",
			wantMsg: "test [Y/n]: ",
			want:    true,
		},
		{
			input:   "hoge\ny",
			msg:     "test",
			wantMsg: "test [Y/n]: Please enter one of the following: Y/y/N/n.\ntest [Y/n]: ",
			want:    true,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ioStream := newTestIOStream()
			ioStream.In.(*bytes.Buffer).WriteString(tt.input)

			cfg := &got.PrompterConfig{
				IsDebug: true,
			}
			p := got.NewPrompter(ioStream, cfg)
			got := p.AskYN(tt.msg)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			stdout := string(ioStream.Out.(*bytes.Buffer).Bytes())
			if diff := cmp.Diff(tt.wantMsg, stdout); diff != "" {
				t.Fatalf("mismatch (-want, +got): %s\n", diff)
			}

			t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
		})
	}
}
