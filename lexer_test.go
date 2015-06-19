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
		"0":         '\000',
		"10":        '\010',
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
			t.Error("bad result for:", sequence)
		}
	}

	for _, str := range []string{"x", "x6", "x6x", "u", "u123", "U123456", "777", "x"} {
		if _, err := readEscapeSequence(bytes.NewBufferString(str)); err == nil {
			t.Error("sequence should cause error:", str)
		}
	}
}

func TestReadLexicalLine(t *testing.T) {
	lines := map[string]LexicalLine{
		"a \\x62 c": LexicalLine{Number: 1, Tokens: []Token{
			{nil, "a", true},
			{nil, "b", false},
			{nil, "c", true},
		}},
		"'a'":     LexicalLine{Number: 1, Tokens: []Token{{nil, "a", false}}},
		"\"b\"":   LexicalLine{Number: 1, Tokens: []Token{{nil, "b", false}}},
		" \"b\" ": LexicalLine{Number: 1, Tokens: []Token{{nil, "b", false}}},
		" \\\" \\' a'b'c'd'": LexicalLine{Number: 1, Tokens: []Token{
			{nil, "\"", false},
			{nil, "'", false},
			{nil, "a'b'c'd'", true},
		}},
		"a (b 'c') d": LexicalLine{Number: 1, Tokens: []Token{
			{nil, "a", true},
			{[]Token{
				{nil, "b", true},
				{nil, "c", false},
			}, "", false},
			{nil, "d", true},
		}},
		"(hey )": LexicalLine{Number: 1, Tokens: []Token{
			{[]Token{{nil, "hey", true}}, "", false},
		}},
		"( \"test\")": LexicalLine{Number: 1, Tokens: []Token{
			{[]Token{{nil, "test", false}}, "", false},
		}},
		"(+ (/ 2 \t3) 4)": LexicalLine{Number: 1, Tokens: []Token{
			{[]Token{
				{nil, "+", true},
				{[]Token{
					{nil, "/", true},
					{nil, "2", true},
					{nil, "3", true},
				}, "", false},
				{nil, "4", true},
			}, "", false},
		}},
		"if a {": LexicalLine{Number: 1, BlockOpen: true, Tokens: []Token{
			{nil, "if", true},
			{nil, "a", true},
		}},
		"'if' a {": LexicalLine{Number: 1, Tokens: []Token{
			{nil, "if", false},
			{nil, "a", true},
			{nil, "{", true},
		}},
		"} catch {": LexicalLine{Number: 1, BlockClose: true, BlockOpen: true, Tokens: []Token{
			{nil, "catch", true},
		}},
		"}": LexicalLine{Number: 1, BlockClose: true, Tokens: []Token{}},
	}
	for str, expected := range lines {
		reader := Lexer{LogicalLineReader{NewPhysLineReader(bytes.NewBufferString(str))}}
		if actual, err := reader.ReadLexicalLine(); err != nil {
			t.Error("got error", err, "for line", str)
		} else if !actual.Equals(&expected) {
			t.Error("got", actual, "but expected", expected, "for line", str)
		}
	}

	errorLines := []string{
		"\"b\"a", "'b'a", "(b)a", "( hey) )", "(hey)'hey'", "(hey)\"hey\"", "'a''b'", "a(hey)",
		"if a{", "if a", "for a", "while a", "try", "else", "def",
		"()",
	}
	for _, line := range errorLines {
		reader := Lexer{LogicalLineReader{NewPhysLineReader(bytes.NewBufferString(line))}}
		if _, err := reader.ReadLexicalLine(); err == nil {
			t.Error("expected error for:", line)
		}
	}
}
