package got

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type PrompterConfig struct {
	IsDebug bool
}

type Prompter struct {
	out io.Writer
	in  io.Reader

	debugL *log.Logger
}

func NewPrompter(ioStream *IOStream, cfg *PrompterConfig) *Prompter {
	return &Prompter{
		out:    ioStream.Out,
		in:     ioStream.In,
		debugL: NewDebugLogger(ioStream.Err, "prompter", cfg.IsDebug),
	}
}

func (p *Prompter) SelectExecutableToDisable(exec *Executable) *Executable {
	p.debugL.Printf("start (*Prompter).SelectExecutableToDisable(%v)\n", exec)

	if exec == nil {
		p.debugL.Printf("end (*Prompter).SelectExecutableToDisable(%v)\n", exec)
		return nil
	}

	fmt.Fprintln(p.out, "Please select an executable to install.")
	exec.Disable = !p.AskYN(exec.Name)

	p.debugL.Printf("end (*Prompter).SelectExecutableToDisable(%v)\n", exec)
	return exec
}

func (p *Prompter) ChooseToForceOverwrite(destPath string) bool {
	p.debugL.Printf("start (*Prompter).ChooseToForceOverwrite(%v)\n", destPath)

	fmt.Fprintf(p.out, "The destination file already exists: %s.\n", destPath)
	fmt.Fprintln(p.out, "Please choose to overwrite or not.")

	isOverwrite := p.AskYN(destPath)

	p.debugL.Printf("end (*Prompter).ChooseToForceOverwrite(%v)\n", destPath)

	return isOverwrite
}

func (p *Prompter) AskYN(msg string) bool {
	p.debugL.Printf("start (*Prompter).AskYN(%s)\n", msg)

	fmt.Fprintf(p.out, "%s [Y/n]: ", msg)
	yn := ""
	fmt.Fscan(p.in, &yn)

	if yn == "" {
		p.debugL.Printf("end (*Prompter).AskYN(%s)\n", msg)
		return true
	}

	if !strings.Contains("YyNn", yn) || len(yn) != 1 {
		fmt.Fprintln(p.out, "Please enter one of the following: Y/y/N/n.")

		p.debugL.Printf("end (*Prompter).AskYN(%s)\n", msg)

		return p.AskYN(msg)
	}

	p.debugL.Printf("end (*Prompter).AskYN(%s)\n", msg)

	return strings.Contains("Yy", yn)
}
