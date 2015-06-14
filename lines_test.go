package pragmash

import (
	"bytes"
	"testing"
)

func TestPhysLineReader(t *testing.T) {
	testPhysLineReaderCase(t, "this\nis\na\ntest", []string{"this", "is", "a", "test"})
	testPhysLineReaderCase(t, "this\nis\r\na\ntest", []string{"this", "is", "a", "test"})
	testPhysLineReaderCase(t, "th\ris\nis\r\na\ntest", []string{"th\ris", "is", "a", "test"})
}

func testPhysLineReaderCase(t *testing.T, str string, lines []string) {
	buffer := bytes.NewBufferString(str)
	reader := NewPhysLineReader(buffer)
	for i, expected := range lines {
		line, num, err := reader.ReadLine()
		if err != nil {
			t.Error(err)
			break
		}
		if num != i+1 {
			t.Error("expected line number", i+1, "but got", num)
		}
		if line != expected {
			t.Error("expected", expected, "but got", line)
		}
	}
}

func TestLogicalLineReader(t *testing.T) {
	testLogLineReaderCase(t, "this\n is\na \t\n\ttest", []string{"this", "is", "a", "test"},
		[]int{1, 2, 3, 4})
	testLogLineReaderCase(t, " this\nis\\\n a\ntest", []string{"this", "is a", "test"},
		[]int{1, 2, 4})
	testLogLineReaderCase(t, "testing\\ \ntesting\\\n123\\\ntesting ", []string{"testing\\",
		"testing123testing"}, []int{1, 2})

	buffer := bytes.NewBufferString("hey\\")
	reader := LogicalLineReader{NewPhysLineReader(buffer)}
	_, _, err := reader.ReadLine()
	if err == nil {
		t.Error("line continuation on last line should cause an error")
	}
}

func testLogLineReaderCase(t *testing.T, str string, lines []string, nums []int) {
	buffer := bytes.NewBufferString(str)
	reader := LogicalLineReader{NewPhysLineReader(buffer)}
	for i, expected := range lines {
		line, num, err := reader.ReadLine()
		if err != nil {
			t.Error(err)
			break
		}
		expectedNum := nums[i]
		if num != expectedNum {
			t.Error("expected line number", expectedNum, "but got", num)
		}
		if line != expected {
			t.Error("expected", expected, "but got", line)
		}
	}
}
