package pragmash

import (
	"errors"
	"strconv"
	"strings"
)

// TokenizeString turns a string into a list of Lines with corresponding
// contexts.
func TokenizeString(str string) ([]Line, []string, error) {
	lines := make([]Line, 0)
	contexts := make([]string, 0)
	tokenizer := NewTokenizer()
	lineStart := -1

	// Loop through each line string
	for i, lineStr := range strings.Split(str, "\n") {
		// If the line is a comment, we should skip it.
		if strings.HasPrefix(strings.TrimSpace(lineStr), "#") {
			continue
		}

		// Add the line to the tokenizer.
		line, err := tokenizer.Line(lineStr)
		if err != nil {
			return nil, nil, err
		} else if line != nil {
			// We read the line, so we should add it and its context.
			lines = append(lines, *line)
			if lineStart >= 0 {
				contexts = append(contexts, "line "+strconv.Itoa(lineStart+1))
				lineStart = -1
			} else {
				contexts = append(contexts, "line "+strconv.Itoa(i+1))
			}
		} else if lineStart < 0 {
			// This line is being continued.
			lineStart = i
		}
	}
	if !tokenizer.Done() {
		return nil, nil, errors.New("unexpected EOF after line continuation")
	}
	return lines, contexts, nil
}

// A Line represents a logical line in a source file.
type Line struct {
	Close  bool
	Open   bool
	Tokens []Token
}

// Blank returns true if the line has no tokens and does not close or open a
// block.
func (l Line) Blank() bool {
	return len(l.Tokens) == 0 && !l.Close && !l.Open
}

// Runnable returns a Runnable for the line.
// If the line is blank, this returns a ValueRunnable containing an empty
// StringValue.
// If the line is not blank, this returns a CommandRunnable.
func (l Line) Runnable(context string) Runnable {
	if len(l.Tokens) == 0 {
		return ValueRunnable{emptyValue}
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
		return ValueRunnable{NewHybridValueString(t.String)}
	}
	if len(t.Nested) == 0 {
		// If a command has no name or arguments, it just returns the empty
		// string.
		return ValueRunnable{NewHybridValueString(t.String)}
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
		t.previous = t.previous + line[0:len(line)-1]
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
		res.Tokens = res.Tokens[0 : len(tokens)-1]
		res.Open = true
	}
	if tokens[0].String == "}" {
		res.Tokens = res.Tokens[1:]
		res.Close = true
	}
	return res, nil
}
