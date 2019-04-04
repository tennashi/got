package got

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Command struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewCommand() *Command {
	return &Command{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (c *Command) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	fmt.Println(cmd)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

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
