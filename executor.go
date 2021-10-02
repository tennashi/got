package got

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandExecutionError struct {
	Command string
	Message string
	Err     error
}

func (e *CommandExecutionError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("failed to execute the command: %s:\n%s: %s", e.Command, e.Message, e.Err.Error())
}

func (e *CommandExecutionError) Unwrap() error {
	return e.Err
}

type ExecutorConfig struct {
	IsDebug bool
}

type Executor struct {
	ioStream *IOStream
	env      map[string][]string

	debugL *log.Logger
}

func NewExecutor(ioStream *IOStream, cfg *ExecutorConfig) *Executor {
	return &Executor{
		ioStream: ioStream,
		env:      make(map[string][]string),
		debugL:   NewDebugLogger(ioStream.Err, "executor", cfg.IsDebug),
	}
}

func (e *Executor) Exec(cmdName string, args []string) error {
	e.debugL.Printf("start (*Executor).Exec(%v, %v)\n", cmdName, args)

	cmd := exec.Command(cmdName, args...)
	e.debugL.Printf("actual command: %s\n", cmd)
	cmd.Stdin = e.ioStream.In
	cmd.Stdout = e.ioStream.Out

	stderrP, err := cmd.StderrPipe()
	if err != nil {
		e.debugL.Printf("error occurred in cmd.StderrPipe(): %s\n", err.Error())

		return err
	}

	stderrR := io.TeeReader(stderrP, e.ioStream.Err)

	env := os.Environ()
	for k, vs := range e.env {
		env = append(env, k+"="+strings.Join(vs, ","))

		e.debugL.Printf("env added: %s=%s\n", k, strings.Join(vs, ","))
	}

	cmd.Env = env

	if err = cmd.Start(); err != nil {
		e.debugL.Printf("error occurred in cmd.Start(): %s\n", err.Error())

		return err
	}

	stderr, err := io.ReadAll(stderrR)
	if err != nil {
		e.debugL.Printf("error occurred in io.ReadAll(): %s\n", err.Error())

		return err
	}

	err = cmd.Wait()
	if err != nil {
		e.debugL.Printf("error occurred in cmd.Wait(): %s\n", err.Error())

		return &CommandExecutionError{
			Command: cmdName + " " + strings.Join(args, " "),
			Message: string(stderr),
			Err:     err,
		}
	}

	e.debugL.Printf("end (*Executor).Exec(%v, %v)\n", cmdName, args)

	return nil
}

func (e *Executor) ExecBackground(cmdName string, args []string) (string, error) {
	e.debugL.Printf("start (*Executor).ExecBackground(%v, %v)\n", cmdName, args)

	cmd := exec.Command(cmdName, args...)
	e.debugL.Printf("actual command: %s\n", cmd)
	cmd.Stdin = e.ioStream.In

	env := os.Environ()
	for k, vs := range e.env {
		env = append(env, k+"="+strings.Join(vs, ","))

		e.debugL.Printf("env added: %s=%s\n", k, strings.Join(vs, ","))
	}

	cmd.Env = env

	stdout, err := cmd.Output()
	if err != nil {
		execErr := &exec.ExitError{}
		if errors.As(err, &execErr) {
			e.debugL.Printf("error occurred in cmd.Output(): %s\n", err.Error())

			return "", &CommandExecutionError{
				Command: cmdName + " " + strings.Join(args, " "),
				Message: string(execErr.Stderr),
				Err:     err,
			}
		}

		return "", err
	}

	e.debugL.Printf("end (*Executor).Exec(%v, %v)\n", cmdName, args)

	return string(stdout), nil
}

func (e *Executor) SetEnv(k, v string) {
	e.debugL.Printf("start (*Executor).SetEnv(%v, %v)\n", k, v)

	e.env[k] = []string{v}

	e.debugL.Printf("end (*Executor).SetEnv(%v, %v)\n", k, v)
}

func (e *Executor) AddEnv(k, v string) {
	e.debugL.Printf("start (*Executor).AddEnv(%v, %v)\n", k, v)

	vs, ok := e.env[k]
	if !ok {
		e.env[k] = []string{v}

		e.debugL.Printf("end (*Executor).AddEnv(%v, %v)\n", k, v)

		return
	}

	e.env[k] = append(vs, v)

	e.debugL.Printf("end (*Executor).AddEnv(%v, %v)\n", k, v)
}
