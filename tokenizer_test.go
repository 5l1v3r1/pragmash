package pragmash

import (
	"testing"
)

func TestForRangeLine(t *testing.T) {
	lineStr := "for x (range 10 0 -1) {"
	lines, _, err := TokenizeString(lineStr)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 {
		t.Fatal("invalid line count")
	}
	
	line := lines[0]
	if !line.Open || line.Close {
		t.Error("line should be open and not close")
	}
	if len(line.Tokens) != 3 {
		t.Fatal("invalid number of tokens")
	}
	
	if line.Tokens[0].String != "for" {
		t.Error("invalid first token")
	}
	if line.Tokens[1].String != "x" {
		t.Error("invalid second token")
	}
	if line.Tokens[2].Nested == nil {
		t.Error("invalid third token")
	} else if len(line.Tokens[2].Nested) != 4 {
		t.Error("invalid third token (len)")
	} else if line.Tokens[2].Nested[0].String != "range" ||
		line.Tokens[2].Nested[1].String != "10" ||
		line.Tokens[2].Nested[2].String != "0" ||
		line.Tokens[2].Nested[3].String != "-1" {
		t.Error("invalid third token (arguments)")
	}
}
