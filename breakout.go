package pragmash

import (
	"errors"
)

// These are the types of breakout which pragmash currently supports.
const (
	BreakoutTypeException = iota
	BreakoutTypeBreak     = iota
	BreakoutTypeContinue  = iota
	BreakoutTypeReturn    = iota
)

// A Breakout is used to jump out of some scope in pragmash.
// Breakouts are used for exceptions, loop control, and return values.
type Breakout struct {
	typeNum int
	context string
	err     error
	value   *Value
}

// NewBreakoutException creates a new exception.
func NewBreakoutException(context string, err error) *Breakout {
	return &Breakout{BreakoutTypeException, context, err, nil}
}

// NewBreakoutBreak creates a new break breakout.
func NewBreakoutBreak(context string) *Breakout {
	return &Breakout{BreakoutTypeBreak, context,
		errors.New("break without loop"), nil}
}

// NewBreakoutContinue creates a new continue breakout.
func NewBreakoutContinue(context string) *Breakout {
	return &Breakout{BreakoutTypeContinue, context,
		errors.New("continue without loop"), nil}
}

// NewBreakoutReturn creates a new return breakout.
func NewBreakoutReturn(context string, val *Value) *Breakout {
	return &Breakout{BreakoutTypeReturn, context,
		errors.New("nothing to return to"), val}
}

// Context returns the context string.
func (b *Breakout) Context() string {
	return b.context
}

// Error returns the error.
func (b *Breakout) Error() error {
	return b.err
}

// Type returns the type of the breakout.
func (b *Breakout) Type() int {
	return b.typeNum
}

// Value returns the value associated with the breakout.
// This is only useful for return breakouts.
func (b *Breakout) Value() *Value {
	return b.value
}
