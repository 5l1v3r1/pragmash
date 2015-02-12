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
func (s Scanner) ReadBare(parenTerm bool) (string, error) {
	// TODO: use bytes.Buffer here to build the string.
	next, _, err := s.ReadRune()
	res := ""
	for err == nil {
		if unicode.IsSpace(next) {
			s.UnreadRune()
			break
		} else if next == '\\' {
			x, err := s.ReadEscape()
			if err != nil {
				return "", err
			}
			res += x
		} else if next == ')' && parenTerm {
			s.UnreadRune()
			break
		} else {
			res += string(next)
		}
		next, _, err = s.ReadRune()
	}
	if err != nil && err != io.EOF {
		return "", err
	}
	return res, nil
}

// ReadCommand reads a command.
// This does not expect to read an open parenthesis as the first character.
func (s Scanner) ReadCommand(parenTerm bool) ([]Token, error) {
	res := []Token{}
	for {
		t, err := s.ReadToken(parenTerm)
		if err != nil {
			return nil, err
		} else if t == nil {
			break
		}
		res = append(res, *t)
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
// This does not expect to read an opening quote as the first character.
func (s Scanner) ReadQuoted() (string, error) {
	// TODO: use bytes.Buffer here
	next, _, err := s.ReadRune()
	res := ""
	for err == nil {
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

// ReadToken reads the next token (i.e. an argument or nested command).
// This will return nil, nil to indicate that the command has no more tokens.
func (s Scanner) ReadToken(parenTerm bool) (*Token, error) {
	s.SkipWhitespace()
	next, _, err := s.ReadRune()
	if err != nil {
		if !parenTerm && err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	if next == '(' {
		args, err := s.ReadCommand(true)
		if err != nil {
			return nil, err
		}
		return &Token{args, ""}, nil
	} else if next == '"' {
		str, err := s.ReadQuoted()
		if err != nil {
			return nil, err
		}
		return &Token{nil, str}, nil
	} else if next == '$' {
		str, err := s.ReadBare(parenTerm)
		if err != nil {
			return nil, err
		}
		return &Token{[]Token{Token{nil, "get"}, Token{nil, str}}, ""}, nil
	} else if next == ')' && parenTerm {
		return nil, nil
	}
	
	s.UnreadRune()
	str, err := s.ReadBare(parenTerm)
	if err != nil {
		return nil, err
	}
	return &Token{nil, str}, nil
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
