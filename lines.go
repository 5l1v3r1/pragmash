package pragmash

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// A LineReader delivers a continual stream of lines until an EOF is
// encountered.
//
// The ReadLine method will return a line and a line number. If an error occurs,
// the error will be non-nil and the other two return values should be ignored.
// If the error is io.EOF, the LineReader has finished reading gracefully.
type LineReader interface {
	ReadLine() (string, int, error)
}

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

// ReadLine reads the next physical line and returns both the line itself and
// the corresponding line number. If an error occurs, it returns a non-nil
// error and the other return values should be ignored.
func (p *PhysLineReader) ReadLine() (string, int, error) {
	result := new(bytes.Buffer)
	for {
		r, _, err := p.reader.ReadRune()
		if err == io.EOF && result.Len() > 0 {
			break
		} else if err != nil {
			return "", 0, err
		} else if r == '\n' {
			break
		}
		if _, err := result.WriteRune(r); err != nil {
			return "", 0, err
		}
	}
	p.lineNumber++
	return strings.TrimSuffix(result.String(), "\r"), p.lineNumber, nil
}

// A LogicalLineReader wraps a LineReader and does extra processing to remove
// whitespace and handle line continuations.
//
// Line numbers from the underlying LineReader are preserved. The first
// underlying line number for the logical line is returned as the logical line
// number.
type LogicalLineReader struct {
	Reader LineReader
}

// ReadLine reads one or more physical lines and returns a logical line or an
// error.
func (l LogicalLineReader) ReadLine() (string, int, error) {
	result := new(bytes.Buffer)
	firstLineNum := -1
	for {
		line, n, err := l.Reader.ReadLine()
		if firstLineNum == -1 {
			firstLineNum = n
		}
		if err == io.EOF && result.Len() > 0 {
			return "", 0, ErrEOFAfterLineContinuation
		} else if err != nil {
			return "", 0, err
		}
		if _, err := result.WriteString(line); err != nil {
			return "", 0, err
		}
		if len(line) == 0 || line[len(line)-1] != '\\' {
			break
		}
	}
	return strings.TrimSpace(result.String()), firstLineNum, nil
}
