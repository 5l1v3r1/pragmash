// Command pragmash executes a pragmash script using the standard library.
package main

import (
	"fmt"
	"github.com/unixpickle/pragmash"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: pragmash <script.pragmash> [ARGS]")
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())

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

	variables := map[string]pragmash.Value{
		"ARGV": pragmash.NewHybridValueString(strings.Join(os.Args[2:], "\n")),
		"DIR": pragmash.NewHybridValueString(filepath.Dir(os.Args[1])),
	}
	runner := pragmash.NewStdRunner(variables)

	if _, exc := runnable.Run(runner); exc != nil {
		fmt.Fprintln(os.Stderr, "exception at "+exc.Context()+": "+
			exc.String())
		os.Exit(1)
	}
}
