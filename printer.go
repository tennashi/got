package got

import (
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
)

type TablePrinterConfig struct {
	IsDebug bool
}

type TablePrinter struct {
	writer *tabwriter.Writer

	debugL *log.Logger
}

func NewTablePrinter(ioStream *IOStream, cfg *TablePrinterConfig) *TablePrinter {
	w := tabwriter.NewWriter(ioStream.Out, 5, 0, 2, ' ', 0)

	return &TablePrinter{
		writer: w,
		debugL: NewDebugLogger(ioStream.Err, "printer", cfg.IsDebug),
	}
}

func (p *TablePrinter) PrintInstalledPackages(pkgs []InstalledPackage) error {
	defer p.writer.Flush()

	fmt.Fprintln(p.writer, "NAME\tVERSION\tINSTALLED EXECUTABLES\t")

	for _, pkg := range pkgs {
		executableNames := make([]string, 0, len(pkg.Executables))
		for _, exec := range pkg.Executables {
			if !exec.Disable {
				executableNames = append(executableNames, exec.Name)
			}
		}

		fmt.Fprintln(p.writer, strings.Join(
			[]string{
				string(pkg.Path),
				pkg.Version,
				strings.Join(executableNames, ","),
			}, "\t"))
	}

	return nil
}
