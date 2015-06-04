package pragmash

import "errors"

var (
	ErrEOFAfterLineContinuation = errors.New("EOF after line continuation")
)
