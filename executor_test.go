package got_test

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tennashi/got"
)

func TestExecutor_Exec(t *testing.T) {
	cfg := &got.ExecutorConfig{
		IsDebug: true,
	}

	t.Run("success", func(t *testing.T) {
		ioStream := newTestIOStream()
		e := got.NewExecutor(ioStream, cfg)

		e.AddEnv("TEST", "test")
		err := e.Exec("sh", []string{"-c", "echo $TEST"})
		if err != nil {
			t.Fatalf("should not be error but: %v", err)
		}

		stdout := string(ioStream.Out.(*bytes.Buffer).Bytes())
		if diff := cmp.Diff("test\n", stdout); diff != "" {
			t.Fatalf("mismatch (-want, +got): %s\n", diff)
		}

		t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
	})

	t.Run("command not found", func(t *testing.T) {
		ioStream := newTestIOStream()
		e := got.NewExecutor(ioStream, cfg)

		err := e.Exec("hoge", nil)
		if err == nil {
			t.Fatalf("should be error but not")
		}

		var execErr *exec.Error
		if !errors.As(err, &execErr) {
			t.Fatalf("want: *exec.Error, but got: %T", err)
		}

		if !strings.Contains(err.Error(), "executable file not found in $PATH") {
			t.Fatalf("unexpected message: %v", err)
		}

		t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
	})

	t.Run("command failure", func(t *testing.T) {
		ioStream := newTestIOStream()
		e := got.NewExecutor(ioStream, cfg)

		err := e.Exec("sh", []string{"-c", "hoge"})
		if err == nil {
			t.Fatalf("should be error but not")
		}

		var execErr *got.CommandExecutionError
		if !errors.As(err, &execErr) {
			t.Fatalf("want: *got.CommandExecutionError, but got: %T", err)
		}

		if !strings.Contains(err.Error(), "hoge: not found") {
			t.Fatalf("unexpected message: %v", err)
		}

		t.Logf("\n%s", ioStream.Err.(*bytes.Buffer).String())
	})

}
