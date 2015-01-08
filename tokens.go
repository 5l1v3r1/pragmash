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
	for reader.Len() > 0 {
		// Consume whitespace.
		for {
			rune, _, _ := reader.ReadRune()
			if !unicode.IsSpace(rune) {
				reader.UnreadRune()
				break
			} else if reader.Len() == 0 {
				return result, nil
			}
		}

		// Read the next argument according to its enclosing context
		// (i.e. whether it's quoted, ticked, or bare).
		next, _, _ := reader.ReadRune()
		var str string
		var err error
		if next == '"' {
			str, err = readString(reader)
		} else if next == '`' {
			str, err = readNestedCommand(reader)
		} else {
			reader.UnreadRune()
			str, err = readBare(reader)
		}

		// Add the argument to the result or fail.
		if err != nil {
			return nil, err
		}
		result = append(result, Token{next == '`', str})
	}

	return result, nil
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

func readNestedCommand(r *strings.Reader) (string, error) {
	// Is this method ugly? Yes. Does it work? Yes. Okay then, glad we settled
	// that.
	var buffer bytes.Buffer
	numOpen := 1
	for r.Len() > 0 {
		next, _, _ := r.ReadRune()
		if next == '`' {
			// If this is the EOL, it should close the nested command.
			if r.Len() == 0 {
				if numOpen != 1 {
					return "", errors.New("Unexpected end of line before ` " +
						"mark.")
				}
				return buffer.String(), nil
			}

			// If the next character is a ` or a space, it's a close tick.
			following, _, _ := r.ReadRune()
			r.UnreadRune()
			isSpace := unicode.IsSpace(following)
			if isSpace || following == '`' {
				numOpen--
				if numOpen == 0 {
					if isSpace {
						return buffer.String(), nil
					} else {
						return "", errors.New("Found excess ` mark in token.")
					}
				}
			} else {
				numOpen++
			}
			buffer.WriteRune(next)
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
	return "", errors.New("Unexpected end of line before `.")
}

func readString(r *strings.Reader) (string, error) {
	var buffer bytes.Buffer
	for r.Len() > 0 {
		next, _, _ := r.ReadRune()
		if next == '"' {
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
