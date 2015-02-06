package pragmash

import (
	"errors"
	"strings"
)

// An Exception stores an error and context info (e.g. a line number) for that
// error.
type Exception struct {
	context string
	err     error
}

// NewException creates a new exception.
func NewException(context string, err error) Exception {
	return Exception{context, err}
}

// Array splits the error string by newlines and returns an array of exceptions,
// each with the same context but different lines from the original error.
func (e Exception) Array() []Value {
	str := e.String()
	if len(str) == 0 {
		return []Value{}
	}
	errorStrs := strings.Split(str, "\n")
	res := make([]Value, len(errorStrs))
	for i, x := range errorStrs {
		res[i] = NewException(e.context, errors.New(x))
	}
	return res
}

// Context returns the context string.
func (e Exception) Context() string {
	return e.context
}

// Error returns the error.
func (e Exception) Error() error {
	return e.err
}

// Number attempts to parse the error's string.
func (e Exception) Number() (Number, error) {
	return ParseNumber(e.String())
}

// String returns the error's string representation (i.e. e.Error().Error())
func (e Exception) String() string {
	return e.err.Error()
}

