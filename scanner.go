package pragmash

import (
	"bytes"
	"errors"
	"io"
	"unicode"
)

// A Scanner reads specific kinds of tokens from a io.RuneScanner.
type Scanner struct {
	io.RuneScanner
}

// NewScannerString creates a scanner which operates on a string.
func NewScannerString(str string) Scanner {
	return Scanner{bytes.NewBufferString(str)}
}

// ReadBare reads a bareword, supporting escapes and terminating at a space or
// an EOF.
func (s Scanner) ReadBare() (string, error) {
	next, _, err := s.ReadRune()
	res := ""
	for err != nil {
		if unicode.IsSpace(next) {
			s.UnreadRune()
			break
		} else if next == '\\' {
			x, err := s.ReadEscape()
			if err != nil {
				return "", err
			}
			res += x
		} else {
			res += string(next)
		}
		next, _, err = s.ReadRune()
	}
	return res, nil
}

// ReadEscape reads escape characters and returns the represented string.
// This does not read a '\'; such a character should already have been read.
func (s Scanner) ReadEscape() (string, error) {
	// TODO: support hex escapes
	next, _, err := s.ReadRune()
	if err != nil {
		return "", errors.New("failed to read escape code")
	}
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

// ReadQuoted reads a quoted string.
// This expects to read an opening quote as the first character.
func (s Scanner) ReadQuoted() (string, error) {
	next, _, err := s.ReadRune()
	if err != nil {
		return "", err
	} else if next != '"' {
		return "", errors.New("Expected to read quotation")
	}
	next, _, err = s.ReadRune()
	res := ""
	for err != nil {
		if next == '"' {
			break
		} else if next == '\\' {
			x, err := s.ReadEscape()
			if err != nil {
				return "", err
			}
			res += x
		} else {
			res += string(next)
		}
		next, _, err = s.ReadRune()
	}
	if err != nil {
		return "", err
	}
	return res, nil
}

// SkipLine reads up to and including the next newline.
func (s Scanner) SkipLine() error {
	next, _, err := s.ReadRune()
	for err == nil {
		if next == '\n' {
			return nil
		}
		next, _, err = s.ReadRune()
	}
	return err
}

// SkipWhitespace reads up to but not including the next non-whitespace
// character.
func (s Scanner) SkipWhitespace() error {
	next, _, err := s.ReadRune()
	for err == nil {
		if !unicode.IsSpace(next) {
			s.UnreadRune()
			return nil
		}
		next, _, err = s.ReadRune()
	}
	return err
}
