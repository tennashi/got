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

func (p *Prompter) SelectExecutableToDisable(pkg *InstalledPackage) {
	p.debugL.Printf("start (*Prompter).SelectExecutableToDisable(%v)\n", pkg)

	if pkg == nil {
		p.debugL.Printf("end (*Prompter).SelectExecutableToDisable(%v)\n", pkg)
		return
	}

	executables := pkg.Executables

	if len(executables) <= 1 {
		p.debugL.Printf("end (*Prompter).SelectExecutableToDisable(%v)\n", pkg)
		return
	}

	fmt.Fprintf(p.out, "Multiple executables are found: %s.\n", pkg.Path)
	fmt.Fprintln(p.out, "Please select an executable to install.")

	for _, executable := range executables {
		executable.Disable = !p.AskYN(executable.Name)
	}

	pkg.Executables = executables

	p.debugL.Printf("end (*Prompter).SelectExecutableToDisable(%v)\n", pkg)
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
