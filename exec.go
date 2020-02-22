package got

import (
	"os/exec"
)

type commandContext struct {
	ioStream *ioStream
	workDir  string
}

func newCommandContext(ioStream *ioStream, dir string) *commandContext {
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
