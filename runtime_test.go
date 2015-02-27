package pragmash

import (
	"strconv"
	"testing"
)

func BenchmarkNumericLoop(b *testing.B) {
	// Generate a script which loops b.N times.
	nString := strconv.Itoa(b.N)
	script := "set x 0\nwhile (< $x " + nString + ") {\n" +
		"set x (+ $x 1)\n}"
	runBenchmarkScript(script)
}

func BenchmarkSummation(b *testing.B) {
	// Generate a script which loops b.N times.
	nString := strconv.Itoa(b.N)
	script := "set x (range " + nString + ")\n" +
		"set sum 0\nfor y $x {\nset sum (+ $sum $y)\n}"
	runBenchmarkScript(script)
}

func runBenchmarkScript(script string) {
	lines, contexts, _ := TokenizeString(script)
	runnable, _ := ScanAll(lines, contexts)
	runner := NewStdRunner(map[string]*Value{})
	runnable.Run(runner)
}
