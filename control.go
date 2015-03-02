package pragmash

import (
	"strings"
)

// A BreakRunner returns a break breakout.
type BreakRunner struct {
	Context string
}

// Run returns nil, NewBreakoutBreak(b.Context).
func (b BreakRunner) Run(r Runner) (*Value, *Breakout) {
	return nil, NewBreakoutBreak(b.Context)
}

// A ContinueRunner returns a continue breakout.
type ContinueRunner struct {
	Context string
}

// Run returns nil, NewBreakoutContinue(c.Context).
func (c ContinueRunner) Run(r Runner) (*Value, *Breakout) {
	return nil, NewBreakoutContinue(c.Context)
}

// A ReturnRunner returns a return breakout.
type ReturnRunner struct {
	Arguments []Runnable
	Context   string
}

// Run generates a return string by running its arguments and joining them.
func (rr ReturnRunner) Run(r Runner) (*Value, *Breakout) {
	// TODO: optimize for single argument returns.
	args := make([]string, len(rr.Arguments))
	for i, x := range rr.Arguments {
		v, bo := x.Run(r)
		if bo != nil {
			return nil, bo
		}
		args[i] = v.String()
	}
	str := NewValueString(strings.Join(args, " "))
	return nil, NewBreakoutReturn(rr.Context, str)
}
