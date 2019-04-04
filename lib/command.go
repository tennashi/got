package got

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Command has stdin/stdout/stderr.
type Command struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewCommand initialize Command with default stdin/stdout/stderr.
func NewCommand() *Command {
	return &Command{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run executes the command.
func (c *Command) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

// RunInDir executes the command in dir.
func (c *Command) RunInDir(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Dir = dir
	return cmd.Run()

}

// SURun executes the command with sudo.
func (c *Command) SURun(command string, args ...string) error {
	if existSudo() {
		execCmd := make([]string, len(args)+1)
		execCmd[0] = command
		for i, arg := range args {
			execCmd[i+1] = arg
		}
		return c.Run("sudo", execCmd...)
	}
	return errors.New("permission denied")
}

func existSudo() bool {
	if _, err := exec.LookPath("sudo"); err != nil {
		return false
	}
	return true
}
