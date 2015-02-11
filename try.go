package pragmash

import (
	"errors"
)

// A Try block represents a try-catch block.
type Try struct {
	Catch        Runnable
	CatchContext string
	Try          Runnable

	// The Variable may be nil if the exception should be discarded.
	Variable Runnable
}

// Run executes the try-catch block.
// This returns an exception if the catch block throws one or if the exception
// variable cannot be set.
func (t Try) Run(r Runner) (Value, *Exception) {
	_, exc := t.Try.Run(r)
	if exc == nil {
		return StringValue(""), nil
	}

	// Set the exception variable if necessary.
	if t.Variable != nil {
		v, e := t.Variable.Run(r)
		if e != nil {
			return nil, e
		}
		if _, e := r.RunCommand("set", []Value{v, *exc}); e != nil {
			return nil, NewException(t.CatchContext, e)
		}
	}

	// Run the catch body
	return t.Catch.Run(r)
}

// A TryScanner scans a try-catch block.
type TryScanner struct {
	catchContext string
	context      string
	scanner      SemanticScanner
	tryBlock     Runnable
	variable     Runnable
}

// NewTryScanner starts a TryScanner or fails if the initiating line is invalid.
func NewTryScanner(l Line, context string) (*TryScanner, error) {
	// Validate the line.
	if len(l.Tokens) != 1 {
		return nil, errors.New("try block takes no arguments")
	} else if l.Tokens[0].String != "try" {
		return nil, errors.New("try block must be initiated by '{'")
	} else if l.Close || !l.Open {
		return nil, errors.New("try line must end with '{' and not start" +
			" with '}'")
	}

	// Generate the result
	return &TryScanner{"", context, newGenericScanner(true), nil, nil}, nil
}

// EOF returns an error with the context of the first line of the loop.
func (t *TryScanner) EOF() (Runnable, error) {
	return nil, errors.New("try block (at " + t.context +
		") not terminated at EOF")
}

// Line adds a line to the try block.
// If the line terminates the block, this returns the block as Runnable.
// If any kind of error is encountered, this returns the error.
// If the block is not closed and the line is properly processed, this returns
// nil, nil.
func (t *TryScanner) Line(l Line, context string) (Runnable, error) {
	if res, err := t.scanner.Line(l, context); err != nil {
		return nil, err
	} else if res != nil {
		if t.tryBlock != nil {
			// The catch block must be closed now.
			if len(l.Tokens) != 0 || l.Open {
				return nil, errors.New("close of catch block (at " + context +
					") takes no arguments")
			}
			return Try{res, t.catchContext, t.tryBlock, t.variable}, nil
		}
		// See if there's no catch block.
		if len(l.Tokens) == 0 && !l.Open {
			return Try{RunnableList{}, context, res, nil}, nil
		}
		// Start reading the catch block
		if len(l.Tokens) != 2 || l.Tokens[0].String != "catch" {
			return nil, errors.New("invalid tokens after try block at " +
				context)
		}
		t.catchContext = context
		t.tryBlock = res
		t.variable = l.Tokens[1].Runnable(context)
		t.scanner = newGenericScanner(true)
	}
	return nil, nil
}
