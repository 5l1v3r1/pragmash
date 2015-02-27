package pragmash

import (
	"errors"
)

// A For block represents a for-loop.
type For struct {
	Body       Runnable
	Context    string
	Expression Runnable

	// The variable is an optional field. If this is nil, the loop has no
	// variable.
	Variable Runnable
}

// Run executes the for loop.
// This fails if the variable name, exression, or body triggers an exception.
func (f For) Run(r Runner) (*Value, *Breakout) {
	expr, bo := f.Expression.Run(r)
	if bo != nil {
		return nil, bo
	}
	var variable *Value
	if f.Variable != nil {
		variable, bo = f.Variable.Run(r)
		if bo != nil {
			return nil, bo
		}
	}
	for _, val := range expr.Array() {
		if variable != nil {
			_, err := r.RunCommand("set", []*Value{variable, val})
			if err != nil {
				return nil, NewBreakoutException(f.Context, err)
			}
		}
		_, bo = f.Body.Run(r)
		if bo == nil || bo.Type() == BreakoutTypeContinue {
			continue
		} else if bo.Type() == BreakoutTypeBreak {
			break
		} else {
			return nil, bo
		}
	}
	return emptyValue, nil
}

// A ForScanner scans a for-loop.
type ForScanner struct {
	context  string
	scanner  SemanticScanner
	value    Runnable
	variable Runnable
}

// NewForScanner starts a ForScanner or fails if the initiating line is invalid.
func NewForScanner(l Line, context string) (*ForScanner, error) {
	// Validate the line.
	if len(l.Tokens) < 2 || len(l.Tokens) > 3 {
		return nil, errors.New("for loop takes one or two arguments")
	} else if l.Tokens[0].String != "for" {
		return nil, errors.New("for loop must start with 'for' token")
	} else if l.Close || !l.Open {
		return nil, errors.New("for line must end with '{' and not start" +
			" with '}'")
	}

	// Generate the result
	res := &ForScanner{context, newGenericScanner(true), nil, nil}
	res.value = l.Tokens[len(l.Tokens)-1].Runnable(context)
	if len(l.Tokens) == 3 {
		res.variable = l.Tokens[1].Runnable(context)
	}
	return res, nil
}

// EOF returns an error with the context of the first line of the loop.
func (f *ForScanner) EOF() (Runnable, error) {
	return nil, errors.New("for loop (at " + f.context +
		") not terminated at EOF")
}

// Line adds a line to the for loop.
// If the line terminates the loop, this returns the loop as Runnable.
// If any kind of error is encountered, this returns the error.
// If the loop is not closed and the line is properly processed, this returns
// nil, nil.
func (f *ForScanner) Line(l Line, context string) (Runnable, error) {
	if res, err := f.scanner.Line(l, context); err != nil {
		return nil, err
	} else if res != nil {
		if len(l.Tokens) > 0 {
			return nil, errors.New("unexpected tokens after for block")
		}
		return For{res, f.context, f.value, f.variable}, nil
	}
	return nil, nil
}
