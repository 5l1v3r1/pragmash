package pragmash

import (
	"io"
	"unicode"
)

func discardLine(r io.RuneScanner) {
	for next, _, err := r.ReadRune(); err == nil; next, _, err = r.ReadRune() {
		if next == '\n' {
			break
		}
	}
}

func discardWhitespace(r io.RuneScanner, allowNewlines bool) {
	for next, _, err := r.ReadRune(); err == nil; next, _, err = r.ReadRune() {
		if !unicode.IsSpace(next) || (!allowNewlines && next == '\n') {
			r.UnreadRune()
			break
		}
	}
}
