package io

import "io"

type Stream struct {
	Out, Err io.Writer
	In       io.Reader
}
