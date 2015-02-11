package pragmash

import (
	"errors"
	"strings"
)

// A BoolValue is a bool which implements the Value interface.
type BoolValue bool

// Array returns an empty array if the receiver is false, or an array with a
// single true BoolValue if it's true.
func (b BoolValue) Array() []Value {
	if !b {
		return []Value{}
	} else {
		return []Value{b}
	}
}

// Context returns an empty string.
func (b BoolValue) Context() string {
	return ""
}

// Number returns an error, since a boolean is not a number.
func (b BoolValue) Number() (Number, error) {
	return nil, errors.New("invalid number: " + b.String())
}

// String returns StringValue("true") for a true receiver and StringValue("")
// for a false one.
func (b BoolValue) String() string {
	if b {
		return "true"
	} else {
		return ""
	}
}

// A StringValue is a string which implements the Value interface.
type StringValue string

// Array splits the string by newline characters and returns a slice of
// StringValues.
func (s StringValue) Array() []Value {
	// Split the string up by newline
	if len(s) == 0 {
		return []Value{}
	}
	parts := strings.Split(string(s), "\n")
	res := make([]Value, len(parts))
	for i, x := range parts {
		res[i] = StringValue(x)
	}
	return res
}

// Context returns an empty string.
func (s StringValue) Context() string {
	return ""
}

// Number attempts to parse the string as a number and returns it.
func (s StringValue) Number() (Number, error) {
	// TODO: perhaps we will cache the numeric result.
	return ParseNumber(string(s))
}

// String casts the receiver to string and returns it.
func (s StringValue) String() string {
	return string(s)
}

// A Value is a read-only variable value.
type Value interface {
	// Array returns the array representation of the value.
	Array() []Value

	// Context returns the context of the value. This is useful if the value is
	// an exception. In most cases, this should be an empty string.
	Context() string

	// Number returns the numerical representation of the value, or an error if
	// the value is not a number.
	Number() (Number, error)

	// String returns the textual representation of the value.
	String() string
}
