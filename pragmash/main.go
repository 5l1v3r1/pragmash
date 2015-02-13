// Command pragmash executes a pragmash script using the standard library.
package main

import (
	"fmt"
	"github.com/unixpickle/pragmash"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: pragmash <script.pragmash> [ARGS]")
		os.Exit(1)
	}
	
	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	lines, contexts, err := pragmash.TokenizeString(string(contents))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	runnable, err := pragmash.ScanAll(lines, contexts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	runner := pragmash.NewStdRunner()
	runner.RunCommand("set", []pragmash.Value{pragmash.StringValue("ARGV"),
		pragmash.StringValue(strings.Join(os.Args[2:], "\n"))})
	
	if _, exc := runnable.Run(runner); exc != nil {
		fmt.Fprintln(os.Stderr, "exception at " + exc.Context() + ": " +
			exc.String())
		os.Exit(1)
	}
}
