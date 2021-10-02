package got_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestPrompter_SelectExecutableToDisable(t *testing.T) {
	cases := []struct {
		input *got.Executable
		stdin string
		want  *got.Executable
	}{
		{
			input: &got.Executable{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: false,
			},
			stdin: "y",
			want: &got.Executable{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: false,
			},
		},
		{
			input: &got.Executable{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: false,
			},
			stdin: "n",
			want: &got.Executable{
				Name:    "got",
				Path:    "/path/to/got",
				Disable: true,
			},
		},
		{
			input: &got.Executable{},
			stdin: "",
			want:  &got.Executable{},
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

func TestPrompter_Select(t *testing.T) {
	cases := []struct {
		input      string
		msg        string
		candidates []string
		wantMsg    string
		want       int
	}{
		{
			input:      "0",
			msg:        "test",
			candidates: []string{"a", "b", "c"},
			wantMsg:    "\ta\n\tb\n\tc\ntest: ",
			want:       0,
		},
		{
			input:      "1",
			msg:        "test",
			candidates: []string{"a", "b", "c"},
			wantMsg:    "\ta\n\tb\n\tc\ntest: ",
			want:       1,
		},
		{
			input:      "2",
			msg:        "test",
			candidates: []string{"a", "b", "c"},
			wantMsg:    "\ta\n\tb\n\tc\ntest: ",
			want:       2,
		},
		{
			input:      "3\n2",
			msg:        "test",
			candidates: []string{"a", "b", "c"},
			wantMsg:    "\ta\n\tb\n\tc\ntest: Invalid input: 3.\n\ta\n\tb\n\tc\ntest: ",
			want:       2,
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
			got := p.Select(tt.msg, tt.candidates, 1)

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
