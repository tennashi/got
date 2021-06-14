package got_test

import (
	"bytes"

	"github.com/tennashi/got"
)

func newTestIOStream() *got.IOStream {
	var testOut bytes.Buffer
	var testErr bytes.Buffer
	var testIn bytes.Buffer

	return &got.IOStream{
		In:  &testIn,
		Out: &testOut,
		Err: &testErr,
	}
}
