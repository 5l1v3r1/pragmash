package pragmash

import (
	"errors"
)

// ScanAll scans every line and context string at once and returns a Runnable or
// fails with an error.
func ScanAll(lines []Line, contexts []string) (Runnable, error) {
	if len(lines) != len(contexts) {
		return nil, errors.New("each line must have exactly one context")
	}
	scanner := NewBodyScanner()
	for i, l := range lines {
		r, err := scanner.Line(l, contexts[i])
		if err != nil {
			return nil, err
		} else if r != nil {
			return nil, errors.New("got premature runnable")
		}
	}
	return scanner.EOF()
}

// A SemanticScanner converts a raw list of Lines into a Runnable.
type SemanticScanner interface {
	// EOF should be called when no more lines will be passed to the scanner.
	// It returns an error if the scanner was expecting more data (i.e. close
	// curly braces), otherwise it returns the block up to the point it was
	// read.
	EOF() (Runnable, error)

	// Line adds a line to the scanner.
	// If the line caused some sort of semantic error, that error is returned.
	// If the line signified the end of the scanner's scope, this returns a
	// Runnable.
	// If the scanner was not done, this returns nil, nil.
	Line(l Line, context string) (Runnable, error)
}

// NewBodyScanner returns a SemanticScanner that will read every line its given
// (provided there are no errors) and return a Runnable (or error) on EOF.
func NewBodyScanner() SemanticScanner {
	return &genericScanner{[]Runnable{}, false, nil, false}
}

// NewSingleScanner returns a SemanticScanner that will read the first complete
// runnable and return it.
// This may consume multiple Lines if the first line it encounters starts a
// block such as a loop or if statement.
func NewSingleScanner() SemanticScanner {
	return &genericScanner{[]Runnable{}, true, nil, false}
}

// A genericScanner reads lines from a file or from inside a scope (i.e. a
// while loop, if statement, etc.)
type genericScanner struct {
	list         []Runnable
	single       bool
	subScanner   SemanticScanner
	waitingClose bool
}

func newGenericScanner(waitingClose bool) *genericScanner {
	return &genericScanner{[]Runnable{}, false, nil, waitingClose}
}

func (g *genericScanner) EOF() (Runnable, error) {
	// If there was a subScanner, we should forward our EOF to it.
	if g.subScanner != nil {
		if res, err := g.subScanner.EOF(); err != nil {
			return nil, err
		} else {
			g.list = append(g.list, res)
			g.subScanner = nil
		}
	}
	if g.waitingClose {
		return nil, errors.New("missing '}' at EOF")
	}
	return RunnableList(g.list), nil
}

func (g *genericScanner) Line(l Line, context string) (Runnable, error) {
	if l.Blank() {
		return nil, nil
	}

	// Pass the line to our sub-scanner if needed.
	if g.subScanner != nil {
		r, err := g.subScanner.Line(l, context)
		if err != nil {
			return nil, err
		} else if r != nil {
			g.subScanner = nil
			g.list = append(g.list, r)
			if g.single {
				return r, nil
			}
		}
		return nil, nil
	}

	// Handle block closes.
	if l.Close {
		if g.waitingClose {
			return RunnableList(g.list), nil
		} else {
			return nil, errors.New("unexpected '}' at " + context)
		}
	}

	// Handle specific constructs.
	if l.Tokens[0].String == "while" {
		whileScanner, err := NewWhileScanner(l, context)
		if err != nil {
			return nil, err
		}
		g.subScanner = whileScanner
		return nil, nil
	} else if l.Tokens[0].String == "for" {
		forScanner, err := NewForScanner(l, context)
		if err != nil {
			return nil, err
		}
		g.subScanner = forScanner
		return nil, nil
	}

	// Handle a regular line
	runnable := l.Runnable(context)
	g.list = append(g.list, runnable)
	if g.single {
		return runnable, nil
	}
	return nil, nil
}
