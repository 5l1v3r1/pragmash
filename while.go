package pragmash

import (
	"errors"
)

// A While loop runs a block of code until a condition is false.
type While struct {
	Body      Runnable
	Condition Runnable
}

// Run runs the while loop.
// On success, this returns an empty string.
func (w While) Run(r Runner) (Value, *Exception) {
	for {
		val, exc := w.Condition.Run(r)
		if exc != nil {
			return nil, exc
		}
		if !val.Bool() {
			break
		}
		_, exc = w.Body.Run(r)
		if exc != nil {
			return nil, exc
		}
	}
	return StringValue(""), nil
}

// A WhileScanner reads a while loop semantically.
type WhileScanner struct {
	condition Runnable
	context   string
	scanner   SemanticScanner
}

// NewWhileScanner starts a WhileScanner or fails if the initiating line is
// invalid.
func NewWhileScanner(l Line, context string) (*WhileScanner, error) {
	// Validate the line's skeleton.
	if !l.Open {
		return nil, errors.New("while loop must open a block")
	} else if l.Close {
		return nil, errors.New("while loop must not close a block")
	} else if len(l.Tokens) == 0 || l.Tokens[0].String != "while" {
		return nil, errors.New("while loop must start with 'while' token")
	}

	// Generate the condition.
	condition := ConditionFromTokens(l.Tokens[1:], context)
	scanner := newGenericScanner(true)
	return &WhileScanner{condition, context, scanner}, nil
}

// EOF returns an error with the context of the first line of the loop.
func (w *WhileScanner) EOF() (Runnable, error) {
	return nil, errors.New("while loop (at " + w.context +
		") not terminated at EOF")
}

// Line adds a line to the while loop.
// If the line terminates the loop, this returns the loop as Runnable.
// If any kind of error is encountered, this returns the error.
// If the loop is not closed and the line is properly processed, this returns
// nil, nil.
func (w *WhileScanner) Line(l Line, context string) (Runnable, error) {
	if res, err := w.scanner.Line(l, context); err != nil {
		return nil, err
	} else if res != nil {
		if len(l.Tokens) > 0 {
			return nil, errors.New("unexpected tokens after while block")
		}
		return While{res, w.condition}, nil
	}
	return nil, nil
}
