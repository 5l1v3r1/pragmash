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

// Runnable returns a Runnable for the line.
// If the line is blank, this returns a ValueRunnable containing an empty
// StringValue.
// If the line is not blank, this returns a CommandRunnable.
func (l Line) Runnable(context string) Runnable {
	if len(l.Tokens) == 0 {
		return ValueRunnable{StringValue("")}
	}

	// Turn arguments into Runnables.
	args := make([]Runnable, len(l.Tokens)-1)
	for i := 1; i < len(l.Tokens); i++ {
		args[i-1] = l.Tokens[i].Runnable(context)
	}

	// Create a CommandRunnable.
	return CommandRunnable{Arguments: args, Context: context,
		Name: l.Tokens[0].Runnable(context)}
}

// A Token is either a raw string or a nested command.
type Token struct {
	Nested []Token
	String string
}

// Runnable returns either a ValueRunnable or a CommandRunnable for the token.
func (t Token) Runnable(context string) Runnable {
	if t.Nested == nil {
		return ValueRunnable{StringValue(t.String)}
	}
	if len(t.Nested) == 0 {
		// If a command has no name or arguments, it just returns the empty
		// string.
		return ValueRunnable{StringValue(t.String)}
	}

	// Turn arguments into Runnables.
	args := make([]Runnable, len(t.Nested)-1)
	for i := 1; i < len(t.Nested); i++ {
		args[i-1] = t.Nested[i].Runnable(context)
	}

	// Create a CommandRunnable.
	return CommandRunnable{Arguments: args, Context: context,
		Name: t.Nested[0].Runnable(context)}
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
		tokens = tokens[0 : len(tokens)-1]
		res.Open = true
	}
	if tokens[0].String == "}" {
		tokens = tokens[1:]
		res.Close = true
	}
	return res, nil
}
