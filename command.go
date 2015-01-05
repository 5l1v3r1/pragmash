package pragmash

import (
	"errors"
	"strings"
)

// Command stores the arguments and name of a command.
type Command struct {
	Name         string
	RawArguments []string
	Arguments    []string
	Raw          string
}

// ParseCommand parses a command and returns it.
// In the case of a syntax error, an error will be returned with details.
// The returned command will have RawArguments, but it will not have Arguments
// until it is flattened.
func ParseCommand(s string) (*Command, error) {
	return nil, errors.New("Not yet implemented.")
}

// Flatten evaluates the sub-commands which this command takes as arguments and
// sets the command's Arguments field.
func (c *Command) Flatten(ctx Context) error {
	c.Arguments = make([]string, len(c.RawArguments))
	for i, arg := range c.Arguments {
		if !strings.HasPrefix(arg, "`") || !strings.HasSuffix(arg, "`") {
			c.Arguments[i] = arg
			continue
		}
		// Evaluate the argument as a command
		cmd := arg[1 : len(arg)-1]
		parsed, err := ParseCommand(cmd)
		if err != nil {
			return err
		}
		if err := parsed.Flatten(ctx); err != nil {
			return err
		}
		value, err := ctx.Evaluate(parsed)
		if err != nil {
			return err
		}
		c.Arguments[i] = value
	}
	return nil
}
