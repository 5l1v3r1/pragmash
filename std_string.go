package pragmash

import (
	"bytes"
	"fmt"
	"regexp"
)

// StdString implements ways of manipulating or creating strings
type StdString struct{}

// Echo joins its arguments with spaces.
func (s StdString) Echo(args []Value) Value {
	interfaceArgs := make([]interface{}, len(args))
	for i, x := range args {
		interfaceArgs[i] = x
	}
	return StringValue(fmt.Sprint(interfaceArgs...))
}

// Join joins its arguments without spaces.
func (s StdString) Join(args []Value) Value {
	var buffer bytes.Buffer
	for _, v := range args {
		buffer.WriteString(v.String())
	}
	return StringValue(buffer.String())
}

// Match runs a regular expression on a string.
func (s StdString) Match(expr, haystack string) (Value, error) {
	// Evaluate the regular expression.
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	res := r.FindAllStringSubmatch(haystack, -1)
	
	// Return the result as a massive list.
	var buffer bytes.Buffer
	for i, x := range res {
		for j, y := range x {
			if i != 0 || j != 0 {
				buffer.WriteRune('\n')
			}
			buffer.WriteString(y)
		}
	}
	return StringValue(buffer.String()), nil
}
