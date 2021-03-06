package pragmash

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"strings"
)

// StdIo implements the standard I/O routines.
type StdIo struct {
	Client http.Client
}

// Cmd executes a shell command and returns its combined output.
func (_ StdIo) Cmd(args ...string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("expected at least 1 argument")
	}

	// Find the command.
	cmdName, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}

	// Run the command.
	cmd := exec.Command(cmdName, args[1:]...)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// Gets reads a line of text from the console.
func (_ StdIo) Gets() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		} else {
			return "", errors.New("end of input")
		}
	}
	return scanner.Text(), nil
}

// HttpCookiesOff disables cookie storage.
func (s StdIo) HttpCookiesOff() {
	s.Client.Jar = nil
}

// HttpCookiesOn enables cookie storage.
func (s StdIo) HttpCookiesOn() {
	// TODO: use a public suffix list
	s.Client.Jar, _ = cookiejar.New(nil)
}

// HttpGet runs an HTTP get request. This respects the current cookie settings.
func (s StdIo) HttpGet(url string, extraHeaders ...string) (string, error) {
	return s.doRequest(url, "GET", nil, extraHeaders)
}

// HttpPost runs an HTTP post request. This respects the current cookie
// settings.
func (s StdIo) HttpPost(url, contentType, body string,
	extraHeaders ...string) (string, error) {
	bodyReader := strings.NewReader(body)
	headers := append([]string{"Content-Type: " + contentType},
		extraHeaders...)
	return s.doRequest(url, "POST", bodyReader, headers)
}

// Print prints text to the console with no newline.
func (_ StdIo) Print(vals ...string) {
	for i, s := range vals {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(s)
	}
}

// Puts prints text to the console with a trailing newline.
func (_ StdIo) Puts(vals ...string) {
	for i, s := range vals {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(s)
	}
	fmt.Println("")
}

// Read reads the contents of a file or a URL.
// This will not use the cookies set by "http" commands.
func (_ StdIo) Read(resource string) (string, error) {
	// Read a web URL if applicable.
	if strings.HasPrefix(resource, "http://") ||
		strings.HasPrefix(resource, "https://") {
		resp, err := http.Get(resource)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	// Read a path.
	contents, err := ioutil.ReadFile(resource)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

// Write writes some data to a file.
func (_ StdIo) Write(path, data string) error {
	return ioutil.WriteFile(path, []byte(data), os.FileMode(0600))
}

// doRequest performs a request given some arguments.
func (s StdIo) doRequest(url, method string, body io.Reader,
	headers []string) (string, error) {
	// Create the request.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	for _, headerLine := range headers {
		idx := strings.Index(headerLine, ": ")
		if idx < 0 {
			return "", errors.New("Invalid header: " + headerLine)
		}
		name := headerLine[:idx]
		value := headerLine[idx+2:]
		req.Header.Set(name, value)
		if strings.ToLower(name) == "host" {
			req.Host = value
		}
	}

	// Send the request.
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
