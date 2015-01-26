package pragmash

import "testing"

func TestTokenize(t *testing.T) {
	testTokenizeCase("a b c", []Token{Token{nil, "a"}, Token{nil, "b"},
		Token{nil, "c"}}, t)
	testTokenizeCase("a \"foo bar\" c", []Token{Token{nil, "a"},
		Token{nil, "foo bar"}, Token{nil, "c"}}, t)
	testTokenizeCase("a foo\\ bar c", []Token{Token{nil, "a"},
		Token{nil, "foo bar"}, Token{nil, "c"}}, t)
	testTokenizeCase("a foo\\\\\\ bar c", []Token{Token{nil, "a"},
		Token{nil, "foo\\ bar"}, Token{nil, "c"}}, t)
	testTokenizeCase("a \\\"hey there\\\" c", []Token{Token{nil, "a"},
		Token{nil, "\"hey"}, Token{nil, "there\""}, Token{nil, "c"}}, t)
	testTokenizeCase("alex (nichol is) cool", []Token{Token{nil, "alex"},
		Token{[]Token{Token{nil, "nichol"}, Token{nil, "is"}}, ""},
		Token{nil, "cool"}}, t)
	testTokenizeCase("alex (nichol (is) cool)", []Token{Token{nil, "alex"},
		Token{[]Token{
			Token{nil, "nichol"},
			Token{[]Token{Token{nil, "is"}}, ""},
			Token{nil, "cool"}}, ""}}, t)
	testTokenizeCase("alex (nichol (is) ) cool", []Token{Token{nil, "alex"},
		Token{[]Token{
			Token{nil, "nichol"},
			Token{[]Token{Token{nil, "is"}}, ""}}, ""},
			Token{nil, "cool"}}, t)
	testTokenizeCase("alex (nichol (is)) cool", []Token{Token{nil, "alex"},
		Token{[]Token{
			Token{nil, "nichol"},
			Token{[]Token{Token{nil, "is"}}, ""}}, ""},
			Token{nil, "cool"}}, t)
	testTokenizeCase("test\\ \\)", []Token{Token{nil, "test )"}}, t)
	testTokenizeCase("\"yo\\\"\"", []Token{Token{nil, "yo\""}}, t)
	testTokenizeError("\"yo", t)
	testTokenizeError("a \"yo hey", t)
	testTokenizeError("\"yo hey", t)
	testTokenizeError("a b\\", t)
	testTokenizeError("a \"b\\\"", t)
	testTokenizeError("\"b\\\"", t)
	testTokenizeError("(foo", t)
	testTokenizeError("(foo (bar)", t)
	testTokenizeError("foo)", t)
	testTokenizeError("foo )", t)
	testTokenizeError("(foo))", t)
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
		if !tokensEqual(toks[i], x) {
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

func tokensEqual(t1 Token, t2 Token) bool {
	if (t1.Command == nil) != (t2.Command == nil) {
		return false
	} else if t1.Command == nil {
		return t1.Text == t2.Text
	} else {
		if len(t1.Command) != len(t2.Command) {
			return false
		}
		for i, x := range t1.Command {
			if !tokensEqual(x, t2.Command[i]) {
				return false
			}
		}
		return true
	}
}
