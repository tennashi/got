package got

import (
	"fmt"
	"io"
	"log"
)

func NewDebugLogger(out io.Writer, scope string, isDebug bool) *log.Logger {
	prefix := fmt.Sprintf("[got][%s][debug] ", scope)
	if !isDebug {
		return log.New(&NopWriter{}, prefix, 0)
	}

	return log.New(out, prefix, 0)
}

func NewLogger(out io.Writer, scope string) *log.Logger {
	return log.New(out, fmt.Sprintf("[got][%s] ", scope), 0)
}
