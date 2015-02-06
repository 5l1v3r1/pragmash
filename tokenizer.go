package pragmash

import (
	"strings"
)

// A Line represents a logical line in a source file.
type Line struct {
	Close  bool
	Open   bool
	Tokens []Token
}

// A Token is either a raw string or a nested command.
type Token struct {
	Nested []Token
	String string
}

// A Tokenizer processes raw lines and returns Lines.
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
// If the line was the end of a command expression, such expression is returned
// as the first return parameter.
//
// If the line needs more data to be completed (i.e. ends with a backslash)
// then this returns nil, nil.
func (t *Tokenizer) Line(line string) (*Line, error) {
	full := t.previous + line
	if strings.HasSuffix(line, "\\") {
		t.previous = full
		return nil, nil
	}
	t.previous = ""
	scanner := NewScannerString(full)
	tokens, err := scanner.ReadCommand(false)
	if err != nil {
		return nil, err
	}
	res := &Line{Tokens: tokens}
	// Check if the line is a close or open block.
	if len(tokens) == 0 {
		return res, nil
	}
	if tokens[len(tokens)-1].String == "{" {
		res.Open = true
	}
	if tokens[0].String == "}" {
		res.Close = true
	}
	return res, nil
}
