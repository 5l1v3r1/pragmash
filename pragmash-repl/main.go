package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/unixpickle/pragmash"
	"math/rand"
	"os"
	"time"
)

const (
	ErrorColor   = "\x1b[91m"
	OutputColor  = "\x1b[2m"
	RegularColor = "\x1b[0m"
)

const (
	RegularPrompt = "> "
	ContPrompt    = "... "
)

func escapeResult(res string) string {
	var buf bytes.Buffer
	buf.WriteRune('"')
	for _, r := range res {
		if r == '\n' {
			buf.WriteString("\\n")
		} else if r == '\\' {
			buf.WriteString("\\\\")
		} else {
			buf.WriteRune(r)
		}
	}
	buf.WriteRune('"')
	return buf.String()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("pragmash version " + pragmash.Version())

	statements := make(chan pragmash.Runnable)
	errorChan := make(chan error)
	conts := make(chan struct{})
	go readInput(statements, errorChan, conts)

	runner := pragmash.NewStdRunner(nil)
	fmt.Print(RegularPrompt)
	for {
		select {
		case statement := <-statements:
			res, err := statement.Run(runner)
			if err != nil {
				fmt.Println(ErrorColor + err.Error().Error() + RegularColor)
			} else {
				fmt.Println(OutputColor + escapeResult(res.String()) +
					RegularColor)
			}
			fmt.Print(RegularPrompt)
		case err := <-errorChan:
			fmt.Println(ErrorColor + err.Error() + RegularColor)
			fmt.Print(RegularPrompt)
		case <-conts:
			fmt.Print(ContPrompt)
		}
	}
}

func readInput(statements chan<- pragmash.Runnable, errorChan chan<- error,
	conts chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	tokenizer := pragmash.NewTokenizer()
	semantic := pragmash.NewSingleScanner()
	for scanner.Scan() {
		line, err := tokenizer.Line(scanner.Text())
		if err != nil {
			errorChan <- err
			tokenizer = pragmash.NewTokenizer()
			semantic = pragmash.NewSingleScanner()
		} else if line != nil {
			stmt, err := semantic.Line(*line, "REPL")
			if err != nil {
				semantic = pragmash.NewSingleScanner()
				errorChan <- err
			} else if stmt != nil {
				statements <- stmt
			} else {
				conts <- struct{}{}
			}
		} else {
			conts <- struct{}{}
		}
	}
	fmt.Println("")
	os.Exit(0)
}
