package pragmash

import "errors"

var (
	ErrEOFAfterLineContinuation = errors.New("EOF after line continuation")
	ErrMissingOpenCurlyBrace    = errors.New("missing open curly brace")
	ErrUnexpectedCloseParen     = errors.New("unexpected ')'")
)
