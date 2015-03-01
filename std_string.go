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
func (_ StdString) Chars(s string) []string {
	runes := []rune(s)
	resArr := make([]string, len(runes))
	for i, x := range runes {
		if x == '\n' {
			resArr[i] = "\\n"
		} else {
			resArr[i] = string(x)
		}
	}
	return resArr
}

// Echo joins its arguments with spaces.
func (_ StdString) Echo(args ...string) string {
	return strings.Join(args, " ")
}

// Escape replaces backslashes with double-backslashes and newlines with "\n".
func (_ StdString) Escape(str string) string {
	s := strings.Replace(str, "\\", "\\\\", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

// Join joins its arguments without spaces.
func (_ StdString) Join(args ...string) string {
	var buffer bytes.Buffer
	for _, s := range args {
		buffer.WriteString(s)
	}
	return buffer.String()
}

// Lowercase joins its arguments with spaces and returns the result, converted
// to lower-case.
func (s StdString) Lowercase(args ...string) string {
	return strings.ToLower(s.Echo(args...))
}

// Match runs a regular expression on a string.
func (_ StdString) Match(expr, haystack string) ([]string, error) {
	// Evaluate the regular expression.
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	res := r.FindAllStringSubmatch(haystack, -1)

	// Return the result as a massive list.
	list := make([]string, 0)
	for _, x := range res {
		for _, y := range x {
			list = append(list, y)
		}
	}
	return list, nil
}

// Rep replaces all occurances of a string with another string.
func (_ StdString) Rep(s, old, replacement string) string {
	return strings.Replace(s, old, replacement, -1)
}

// Substr returns a substring of a large string.
func (_ StdString) Substr(s string, start, end int) string {
	if len(s) == 0 {
		return ""
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

	return s[start:end]
}

// Unescape replaces "\\" with "\" and "\n" with a newline.
func (_ StdString) Unescape(arg string) string {
	s := strings.Replace(arg, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

// Uppercase joins its arguments with spaces and returns the result, converted
// to upper-case.
func (s StdString) Uppercase(args ...string) string {
	return strings.ToUpper(s.Echo(args...))
}
