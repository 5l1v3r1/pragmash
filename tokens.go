package pragmash

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

// Token represents a string or an embedded command.
type Token struct {
	Command bool
	Text    string
}

// Tokenize intelligently splits a line into tokens.
// Whitespace will be ignored and escapes will be accounted for.
func Tokenize(line string) ([]Token, error) {
	reader := strings.NewReader(line)
	result := make([]Token, 0)

	for {
		if token, err := readArgument(reader); err != nil {
			return nil, err
		} else if token == nil {
			break
		} else {
			result = append(result, *token)
		}
	}

	return result, nil
}

func readArgument(r *strings.Reader) (*Token, error) {
	discardWhitespace(r, false)

	next, _, err := r.ReadRune()
	if err != nil {
		return nil, nil
	}

	var res Token
	if next == '"' {
		res.Text, err = readString(r)
	} else if next == '`' {
		res.Text, err = readNestedCommand(r)
		res.Command = true
	} else if next == '$' {
		res.Text, err = readBare(r)
		res.Text = "get " + res.Text
		res.Command = true
	} else {
		r.UnreadRune()
		res.Text, err = readBare(r)
	}

	// Add the argument to the result or fail.
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func readBare(r *strings.Reader) (string, error) {
	var buffer bytes.Buffer
	for r.Len() > 0 {
		next, _, _ := r.ReadRune()
		if next == '\\' {
			str, err := readEscape(r)
			if err != nil {
				return "", err
			}
			buffer.WriteString(str)
		} else if unicode.IsSpace(next) {
			break
		} else {
			buffer.WriteRune(next)
		}
	}
	return buffer.String(), nil
}

func readEscape(r *strings.Reader) (string, error) {
	// TODO: support hex escapes
	if r.Len() == 0 {
		return "", errors.New("Cannot escape end of file.")
	}
	next, _, _ := r.ReadRune()
	if next == 'n' {
		return "\n", nil
	} else if next == 'r' {
		return "\r", nil
	} else if next == 'a' {
		return "\a", nil
	} else if next == 't' {
		return "\t", nil
	}
	return string(next), nil
}

func readNestedCommand(r *strings.Reader) (string, error) {
	var buffer bytes.Buffer
	closed := false
	for r.Len() > 0 {
		next, _, _ := r.ReadRune()
		if next == '`' {
			if r.Len() != 0 {
				following, _, _ := r.ReadRune()
				r.UnreadRune()
				if !unicode.IsSpace(following) {
					return "", errors.New("Unexpected character following " +
						"close-tick: '" + string(following) + "'")
				}
			}
			closed = true
			break
		} else if next == '\\' {
			str, err := readEscape(r)
			if err != nil {
				return "", err
			}
			buffer.WriteString(str)
		} else {
			buffer.WriteRune(next)
		}
	}
	if !closed {
		return "", errors.New("Unexpected end of line before `.")
	}
	return buffer.String(), nil
}

func readString(r *strings.Reader) (string, error) {
	var buffer bytes.Buffer
	foundQuote := false
	for r.Len() > 0 {
		next, _, _ := r.ReadRune()
		if next == '"' {
			foundQuote = true
			break
		} else if next == '\\' {
			str, err := readEscape(r)
			if err != nil {
				return "", err
			}
			buffer.WriteString(str)
		} else {
			buffer.WriteRune(next)
		}
	}
	if !foundQuote {
		return "", errors.New("Missing expected \" at end of line.")
	}
	return buffer.String(), nil
}
