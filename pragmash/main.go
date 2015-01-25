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
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "<script> [ARGS]")
		os.Exit(1)
	}
	
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	program, err := pragmash.ParseProgram(string(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	context := pragmash.NewStandardContext()
	context.Variables["ARGV"] = strings.Join(os.Args[2:], "\n")
	if _, err := program.Run(context); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
