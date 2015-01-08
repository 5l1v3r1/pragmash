package pragmash

import "testing"

func TestParseScript(t *testing.T) {
	testParseScriptCase(t, "hey\nthere\ntesting", "basic1", []string{"hey",
		"there", "testing"}, []int{0, 1, 2}, []int{1, 1, 1})
	testParseScriptCase(t, "hey\nthere\ntesting\n", "basic2", []string{"hey",
		"there", "testing"}, []int{0, 1, 2}, []int{1, 1, 1})
	testParseScriptCase(t, "hey\\\nthere\ntesting", "onewrap1",
		[]string{"heythere", "testing"}, []int{0, 2}, []int{2, 1})
	testParseScriptCase(t, "hey\\\nthere\ntesting\n", "onewrap2",
		[]string{"heythere", "testing"}, []int{0, 2}, []int{2, 1})
	testParseScriptCase(t, "foo\nhey\\\nthere\\\ntesting", "twowrap1",
		[]string{"foo", "heytheretesting"}, []int{0, 1}, []int{1, 3})
	testParseScriptCase(t, "foo\nhey\\\nthere\\\ntesting\n", "twowrap2",
		[]string{"foo", "heytheretesting"}, []int{0, 1}, []int{1, 3})
}

func testParseScriptCase(t *testing.T, script, scriptId string, lines []string,
	starts, lens []int) {
	s, err := ParseScript(script)
	if err != nil {
		t.Error(err)
		return
	}
	if s.Len() != len(lines) || len(s.LineStarts) != len(starts) ||
		 len(s.LineLens) != len(lens) {
		t.Error("Unexpected result for script", scriptId)
		return
	}
	for i, x := range s.LogicalLines {
		if x != lines[i] || s.LineStarts[i] != starts[i] ||
			s.LineLens[i] != lens[i] {
			t.Error("Unexpected result for script", scriptId)
			return
		}
	}
}
