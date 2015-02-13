package pragmash

import (
	"bytes"
	"regexp"
	"strings"
)

// StdString implements ways of manipulating or creating strings
type StdString struct{}

// Echo joins its arguments with spaces.
func (s StdString) Echo(args []Value) Value {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = x.String()
	}
	return StringValue(strings.Join(strArgs, " "))
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

// Rep replaces all occurances of a string with another string.
func (_ StdString) Rep(s, old, replacement string) Value {
	return StringValue(strings.Replace(s, old, replacement, -1))
}

// Substr returns a substring of a large string.
func (_ StdString) Substr(s string, start, end int) Value {
	if len(s) == 0 {
		return StringValue("")
	}

	// Any inputs are sanitized and accepted.
	if start < 0 {
		start = 0
	} else if start > len(s) {
		start = len(s)
	}
	if end < start {
		end = start
	} else if end > len(s) {
		end = len(s)
	}

	return StringValue(s[start:end])
}
