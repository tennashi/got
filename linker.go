package got

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type AlreadyExistsError struct {
	Path string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("there is the file in the path: %s", e.Path)
}

type ExecutableLinkerConfig struct {
	BinDir  string
	IsDebug bool
}

type ExecutableLinker struct {
	binDir string

	debugL *log.Logger
}

func NewExecutableLinker(ioStream *IOStream, cfg *ExecutableLinkerConfig) (*ExecutableLinker, error) {
	if cfg.BinDir == "" {
		return nil, &InvalidParamError{Message: "bin directory must not be empty"}
	}

	return &ExecutableLinker{
		binDir: cfg.BinDir,
		debugL: NewDebugLogger(ioStream.Err, "linker", cfg.IsDebug),
	}, nil
}

func (l *ExecutableLinker) Link(executable *Executable) error {
	l.debugL.Printf("start (*ExecutableLinker).Link(%v)\n", executable)
	if executable.Disable {
		l.debugL.Printf("skip symlink because it is disabled: %s\n", executable.Name)
		return nil
	}

	destPath := filepath.Join(l.binDir, executable.Name)
	err := os.Symlink(executable.Path, destPath)
	if err != nil {
		if os.IsExist(err) {
			return &AlreadyExistsError{Path: destPath}
		}
		return err
	}

	l.debugL.Printf("end (*ExecutableLinker).Link(%v)\n", executable)
	return nil
}

func (l *ExecutableLinker) ForceLink(executable *Executable) error {
	l.debugL.Printf("start (*ExecutableLinker).ForceLink(%v)\n", executable)

	if executable.Disable {
		l.debugL.Printf("skip symlink because it is disabled: %s\n", executable.Name)
		return nil
	}

	destPath := filepath.Join(l.binDir, executable.Name)

	if err := os.Remove(destPath); err != nil {
		return err
	}

	if err := os.Symlink(executable.Path, destPath); err != nil {
		return err
	}

	l.debugL.Printf("end (*ExecutableLinker).ForceLink(%v)\n", executable)
	return nil
}
