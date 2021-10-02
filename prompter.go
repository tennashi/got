package got

import (
	"fmt"
	"io"
	"log"
	"strconv"
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

func (p *Prompter) SelectPackage(pkgs []InstalledPackage) *InstalledPackage {
	p.debugL.Printf("start (*Prompter).SelectPackage(%v)\n", pkgs)

	pkgNames := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		pkgNames = append(pkgNames, string(pkg.Path))
	}

	pkg := pkgs[p.Select("Please select a package", pkgNames, 0)]

	p.debugL.Printf("end (*Prompter).SelectPackage(%v)\n", pkgs)

	return &pkg
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

func (p *Prompter) Select(msg string, candidates []string, defaultIndex int) int {
	p.debugL.Printf("start (*Prompter).Select(%s, %v, %d)\n", msg, candidates, defaultIndex)

	if defaultIndex >= len(candidates) {
		defaultIndex = len(candidates) - 1
	}

	for _, candidate := range candidates {
		fmt.Fprintf(p.out, "\t%s\n", candidate)
	}

	fmt.Fprintf(p.out, "%s: ", msg)

	input := ""
	fmt.Fscan(p.in, &input)

	if input == "" {
		p.debugL.Printf("end (*Prompter).Select(%s, %v, %d)\n", msg, candidates, defaultIndex)

		return defaultIndex
	}

	index, err := strconv.Atoi(input)
	if err != nil || index >= len(candidates) {
		fmt.Fprintf(p.out, "Invalid input: %s.\n", input)

		p.debugL.Printf("end (*Prompter).Select(%s, %v, %d)\n", msg, candidates, defaultIndex)

		return p.Select(msg, candidates, defaultIndex)
	}

	p.debugL.Printf("end (*Prompter).Select(%s, %v, %d)\n", msg, candidates, defaultIndex)

	return index
}
