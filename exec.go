package got

import (
	"os/exec"

	"github.com/tennashi/got/io"
)

type commandContext struct {
	ioStream *io.Stream
	workDir  string
}

func newCommandContext(ioStream *io.Stream, dir string) *commandContext {
	return &commandContext{
		ioStream: ioStream,
		workDir:  dir,
	}
}

func (c *commandContext) Exec(cmdName string, args []string) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = c.ioStream.In
	cmd.Stdout = c.ioStream.Out
	cmd.Stderr = c.ioStream.Err
	cmd.Dir = c.workDir
	return cmd.Run()
}
