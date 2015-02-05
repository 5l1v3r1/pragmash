package pragmash

import (
	"strings"
)

type Line struct {
	CloseBlock bool
	OpenBlock  bool
	Tokens     []string
}

type Tokenizer struct {
	previous string
}

// NewTokenizer returns a tokenizer with an empty buffer.
func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

// Done returns true if the tokenizer has no text waiting to be processed.
func (t *Tokenizer) Done() bool {
	return len(t.previous) == 0
}

// Line takes a line as a string and processes it.
//
// If the line was the end of a logical token expression, such expression is
// returned as the first return parameter.
//
// If the line needs more data to be completed (i.e. ends with a backslash)
// then this returns a nil *Line and a nil error.
func (t *Tokenizer) Line(line string) (*Line, error) {
	full := t.previous + line
	if strings.HasSuffix(line, "\\") {
		t.previous = full
		return nil, nil
	}
	// TODO: tokenize the line here
	return nil, nil
}
