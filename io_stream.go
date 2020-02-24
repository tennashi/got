package got

import (
	"bytes"
	"io"
	"os"
)

type ioStream struct {
	Out, Err io.Writer
	In       io.Reader
}

var defaultIOStream = &ioStream{
	In:  os.Stdin,
	Out: os.Stdout,
	Err: os.Stderr,
}

var testOut bytes.Buffer
var testErr bytes.Buffer
var testIn bytes.Buffer

var testIOStream = &ioStream{
	In:  &testIn,
	Out: &testOut,
	Err: &testErr,
}
