package pragmash

import "testing"

func TestTokenize(t *testing.T) {
	testTokenizeCase("a b c", []Token{Token{false, "a"}, Token{false, "b"},
		Token{false, "c"}}, t)
	testTokenizeCase("a \"foo bar\" c", []Token{Token{false, "a"},
		Token{false, "foo bar"}, Token{false, "c"}}, t)
	testTokenizeCase("a foo\\ bar c", []Token{Token{false, "a"},
		Token{false, "foo bar"}, Token{false, "c"}}, t)
	testTokenizeCase("a foo\\\\\\ bar c", []Token{Token{false, "a"},
		Token{false, "foo\\ bar"}, Token{false, "c"}}, t)
	testTokenizeCase("a \\\"hey there\\\" c", []Token{Token{false, "a"},
		Token{false, "\"hey"}, Token{false, "there\""}, Token{false, "c"}}, t)
	testTokenizeCase("alex `nichol is` cool", []Token{Token{false, "alex"},
		Token{true, "nichol is"}, Token{false, "cool"}}, t)
	testTokenizeCase("alex `nichol \\`is\\`` cool", []Token{Token{false, "alex"},
		Token{true, "nichol `is`"}, Token{false, "cool"}}, t)
	testTokenizeCase("alex `nichol \\`is\\` ` cool", []Token{Token{false, "alex"},
		Token{true, "nichol `is` "}, Token{false, "cool"}}, t)
	testTokenizeCase("`test\\``", []Token{Token{true, "test`"}}, t)
	testTokenizeCase("\"yo\\\"\"", []Token{Token{false, "yo\""}}, t)
	testTokenizeError("\"yo", t)
	testTokenizeError("a \"yo hey", t)
	testTokenizeError("\"yo hey", t)
	testTokenizeError("a b\\", t)
	testTokenizeError("a \"b\\\"", t)
	testTokenizeError("\"b\\\"", t)
	testTokenizeError("`foo", t)
	testTokenizeError("`foo``bar`", t)
}

func testTokenizeCase(raw string, toks []Token, t *testing.T) {
	parsed, err := Tokenize(raw)
	if err != nil {
		t.Error(err)
		return
	}
	if len(parsed) != len(toks) {
		t.Error("Bad result for:", raw, "got", parsed)
		return
	}
	for i, x := range parsed {
		tok := toks[i]
		if tok.Command != x.Command || tok.Text != x.Text {
			t.Error("Bad result for:", raw, "got", parsed)
			return
		}
	}
}

func testTokenizeError(raw string, t *testing.T) {
	_, err := Tokenize(raw)
	if err == nil {
		t.Error("Expected error for:", raw)
	}
}
