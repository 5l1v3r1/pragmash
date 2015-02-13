package pragmash

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// StdIo implements the standard I/O routines.
type StdIo struct{}

// Gets reads a line of text from the console.
func (s StdIo) Gets() (Value, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		} else {
			return nil, errors.New("end of input")
		}
	}
	return StringValue(scanner.Text()), nil
}

// Print prints text to the console with no newline.
func (s StdIo) Print(vals []Value) {
	for i, v := range vals {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(v.String())
	}
}

// Puts prints text to the console with a trailing newline.
func (s StdIo) Puts(vals []Value) {
	for i, v := range vals {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(v.String())
	}
	fmt.Println("")
}

// Read reads the contents of a file or a URL.
func (s StdIo) Read(resource string) (Value, error) {
	// Read a web URL if applicable.
	if strings.HasPrefix(resource, "http://") ||
		strings.HasPrefix(resource, "https://") {
		resp, err := http.Get(resource)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return StringValue(body), nil
	}

	// Read a path.
	contents, err := ioutil.ReadFile(resource)
	if err != nil {
		return nil, err
	}
	return StringValue(contents), nil
}

// Write writes some data to a file.
func (s StdIo) Write(path, data string) error {
	return ioutil.WriteFile(path, []byte(data), os.FileMode(0600))
}
