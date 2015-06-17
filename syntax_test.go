package pragmash

import (
	"bytes"
	"testing"
)

func TestReadEscapeSequence(t *testing.T) {
	codes := map[string]rune{
		"a":         '\a',
		"b":         '\b',
		"f":         '\f',
		"n":         '\n',
		"r":         '\r',
		"t":         '\t',
		"v":         '\v',
		"x6a":       'j',
		"x6A":       'j',
		"123":       'S',
		"u2702":     '\u2702',
		"U0001F601": '\U0001F601',
	}
	for _, str := range []string{"(", ")", "?", "'", "\"", "\\", " "} {
		codes[str] = []rune(str)[0]
	}
	for sequence, expected := range codes {
		actual, err := readEscapeSequence(bytes.NewBufferString(sequence))
		if err != nil {
			t.Error(err)
		} else if actual != expected {
			t.Error("bad result for", sequence)
		}
	}
}
