package pragmash

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// A Script maps raw lines in a file to logical lines in a script.
// In source files, a backslash may be added before newlines to indicate that
// the "line" continues.
//
// For example, the following file has two lines but only one logical line:
//
//     cat /foo/bar \
//       /bar/foo`
//
// In this case, the Script would be Script{[]string{"cat /foo/bar /bar/foo"},
// []int{0}, []int{2}}.
type Script struct {
	LogicalLines []string
	LineStarts   []int
	LineLens     []int
}

// ParseScript reads the logical lines of a raw script and returns a Script
// instance.
// This automatically removes empty lines, comments, and whitespace from the
// beginning of each logical line.
func ParseScript(script string) (*Script, error) {
	buffer := strings.NewReader(script)
	realLine := 0
	res := Script{[]string{}, []int{}, []int{}}
	for buffer.Len() > 0 {
		line, count, err := readLogicalLine(buffer)
		if err != nil {
			return nil, err
		}
		
		// Add the line to the result if it's not blank.
		if len(line) > 0 {
			res.LogicalLines = append(res.LogicalLines, line)
			res.LineStarts = append(res.LineStarts, realLine)
			res.LineLens = append(res.LineLens, count)
		}
		
		realLine += count
	}
	return &res, nil
}

// Get returns the logical line at a given index.
func (s *Script) Get(idx int) string {
	return s.LogicalLines[idx]
}

// Len returns the number of logical lines in the script.
func (s *Script) Len() int {
	return len(s.LogicalLines)
}

// Range returns the line number and number of lines in the original file
// corresponding to the specified logical line.
func (s *Script) Range(idx int) (start int, len int) {
	return s.LineStarts[idx], s.LineLens[idx]
}

func readLogicalLine(r io.RuneScanner) (string, int, error) {
	count := 1
	res := bytes.Buffer{}
	discardWhitespace(r, false)
	// Note: we must respect escaped newlines and comments.
	for next, _, err := r.ReadRune(); err == nil; next, _, err = r.ReadRune() {
		if next == '\n' {
			break
		} else if res.Len() == 0 && next == '#' {
			discardLine(r)
			break
		} else if next == '\\' {
			following, _, anErr := r.ReadRune()
			if anErr != nil {
				return "", -1, errors.New("Unexpected backslash at end of " +
					"script")
			}
			if following == '\n' {
				count++
				continue
			}
			r.UnreadRune()
		}
		res.WriteRune(next)
	}
	
	return res.String(), count, nil
}
