package pragmash

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// A PhysLineReader wraps an io.Reader and reads physical lines from it.
type PhysLineReader struct {
	reader     *bufio.Reader
	lineNumber int
}

// NewPhysLineReader creates a PhysLineReader which acts on a given io.Reader
// and starts at line 1.
func NewPhysLineReader(r io.Reader) *PhysLineReader {
	return &PhysLineReader{bufio.NewReader(r), 0}
}

// ReadPhysicalLine reads the next physical line and returns both the line
// itself and the corresponding line number. If an error occurs, it returns a
// non-nil error and the other return values should be ignored.
func (p *PhysLineReader) ReadPhysicalLine() (string, int, error) {
	result := new(bytes.Buffer)
	for {
		r, _, err := p.reader.ReadRune()
		if err == io.EOF || r == '\n' {
			break
		} else if err != nil {
			return "", 0, err
		}
		result.WriteRune(r)
	}
	p.lineNumber++
	return strings.TrimSuffix(result.String(), "\r"), p.lineNumber, nil
}
