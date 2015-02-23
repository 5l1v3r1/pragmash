package pragmash

import (
	"bytes"
	"regexp"
	"strings"
)

// StdString implements ways of manipulating or creating strings
type StdString struct{}

// Chars returns an array with each character from a string.
// Each newline character will be encoded as the two-character escape
// sequence "\\n".
func (_ StdString) Chars(s string) Value {
	runes := []rune(s)
	resArr := make([]Value, len(runes))
	for i, x := range runes {
		if x == '\n' {
			resArr[i] = NewHybridValueString("\\n")
		} else {
			resArr[i] = NewHybridValueString(string(x))
		}
	}
	return NewHybridValueArray(resArr)
}

// Echo joins its arguments with spaces.
func (_ StdString) Echo(args []Value) Value {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = x.String()
	}
	return NewHybridValueString(strings.Join(strArgs, " "))
}

// Join joins its arguments without spaces.
func (_ StdString) Join(args []Value) Value {
	var buffer bytes.Buffer
	for _, v := range args {
		buffer.WriteString(v.String())
	}
	return NewHybridValueString(buffer.String())
}

// Lowercase joins its arguments with spaces and returns the result, converted
// to lower-case.
func (_ StdString) Lowercase(args []Value) Value {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = strings.ToLower(x.String())
	}
	return NewHybridValueString(strings.Join(strArgs, " "))
}

// Match runs a regular expression on a string.
func (_ StdString) Match(expr, haystack string) (Value, error) {
	// Evaluate the regular expression.
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	res := r.FindAllStringSubmatch(haystack, -1)

	// Return the result as a massive list.
	list := make([]Value, 0)
	for _, x := range res {
		for _, y := range x {
			list = append(list, NewHybridValueString(y))
		}
	}
	return NewHybridValueArray(list), nil
}

// Rep replaces all occurances of a string with another string.
func (_ StdString) Rep(s, old, replacement string) Value {
	return NewHybridValueString(strings.Replace(s, old, replacement, -1))
}

// Substr returns a substring of a large string.
func (_ StdString) Substr(s string, start, end int) Value {
	if len(s) == 0 {
		return emptyValue
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

	return NewHybridValueString(s[start:end])
}

// Uppercase joins its arguments with spaces and returns the result, converted
// to upper-case.
func (_ StdString) Uppercase(args []Value) Value {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = strings.ToUpper(x.String())
	}
	return NewHybridValueString(strings.Join(strArgs, " "))
}
