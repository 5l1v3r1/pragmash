package pragmash

import (
	"bytes"
	"errors"
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
func ParseScript(script string) (*Script, error) {
	buffer := strings.NewReader(script)
	realLine := 0
	res := Script{[]string{}, []int{}, []int{}}
	for buffer.Len() > 0 {
		// Read the next logical line.
		logLine := bytes.Buffer{}
		start := realLine
		for buffer.Len() > 0 {
			next, _, _ := buffer.ReadRune()
			if next == '\\' {
				if buffer.Len() == 0 {
					return nil, errors.New("Unexpected backslash at end of " +
						"script")
				}
				if following, _, _ := buffer.ReadRune(); following == '\n' {
					realLine++
					continue
				}
				buffer.UnreadRune()
			} else if next == '\n' {
				break
			}
			logLine.WriteRune(next)
		}
		realLine++
		res.LogicalLines = append(res.LogicalLines, logLine.String())
		res.LineStarts = append(res.LineStarts, start)
		res.LineLens = append(res.LineLens, realLine-start)
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
