package got

import (
	"io"
	"os"
)

type IOStream struct {
	Out, Err io.Writer
	In       io.Reader
}

var defaultIOStream = &IOStream{
	In:  os.Stdin,
	Out: os.Stdout,
	Err: os.Stderr,
}
