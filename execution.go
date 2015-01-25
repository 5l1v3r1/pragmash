package pragmash

import (
	"strings"
)

// Argument is an argument for a command.
// An argument is either a raw string or a sub-command.
type Argument struct {
	Text    string
	Command *Command
}

// Run evaluates an argument and returns it.
// If the argument is plain text, this will never return an error.
func (a *Argument) Run(c Context) (string, error) {
	if a.Command == nil {
		return a.Text, nil
	}
	return a.Command.Run(c)
}

// Block stores a command or a control-flow block.
// A block can be executed in any context.
type Block interface {
	Run(c Context) (string, error)
}

// Blocks is zero or more blocks which execute in order.
type Blocks []Block

// Run executes each block in order and returns the first error it encounters.
// Upon success, this returns an empty string and a nil error.
func (b Blocks) Run(c Context) (string, error) {
	for _, x := range b {
		_, err := x.Run(c)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// Command is a block representing a single-line command.
type Command struct {
	Name      Argument
	Arguments []Argument
}

// Run executes the command.
func (c *Command) Run(ctx Context) (string, error) {
	name, err := c.Name.Run(ctx)
	if err != nil {
		return "", err
	}
	args := make([]string, len(c.Arguments))
	for i, x := range c.Arguments {
		val, err := x.Run(ctx)
		if err != nil {
			return "", err
		}
		args[i] = val
	}
	return ctx.Run(name, args)
}

// Condition is a condition used for if-statements.
type Condition []Argument

// Evaluate runs the condition and returns true if the condition is true.
// An error is returned if any commands in the condition failed to run.
func (c Condition) Evaluate(ctx Context) (bool, error) {
	if len(c) == 0 {
		return true, nil
	}

	// We will need the first argument regardless.
	val, err := c[0].Run(ctx)
	if err != nil {
		return false, err
	}

	// Single argument conditional uses stringy booleans.
	if len(c) == 1 {
		return val != "", nil
	}

	// Make sure the rest of the arguments equal the first one.
	for i := 1; i < len(c); i++ {
		aVal, err := c[i].Run(ctx)
		if err != nil {
			return false, err
		}
		if aVal != val {
			return false, nil
		}
	}
	return true, nil
}

// A Context is used to execute Blocks.
type Context interface {
	Run(command string, args []string) (string, error)
}

// ForBlock is a "for" loop.
type ForBlock struct {
	Variable   *Argument
	Expression Argument
	Block      Block
}

// Run executes the "for" loop.
func (f *ForBlock) Run(c Context) (string, error) {
	body, err := f.Expression.Run(c)
	if err != nil {
		return "", err
	} else if len(body) == 0 {
		return "", nil
	}

	if f.Variable == nil {
		count := strings.Count(body, "\n") + 1
		for i := 0; i < count; i++ {
			if _, err := f.Block.Run(c); err != nil {
				return "", err
			}
		}
		return "", nil
	}

	varName, err := f.Variable.Run(c)
	if err != nil {
		return "", err
	}

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if _, err := c.Run("set", []string{varName, line}); err != nil {
			return "", err
		}
		if _, err := f.Block.Run(c); err != nil {
			return "", err
		}
	}

	return "", nil
}

// IfBlock is an "if" statement.
type IfBlock struct {
	// An array of conditions.
	Conditions []Condition

	// An array of code blocks.
	// Except for the last block, which may be the "else" section, each element
	// in this array should correspond to a condition. Thus, len(Branches) must
	// be either len(Conditions) or len(Conditions)+1.
	Branches []Block
}

// Run executes the if statement.
func (i *IfBlock) Run(c Context) (string, error) {
	for j, x := range i.Conditions {
		val, err := x.Evaluate(c)
		if err != nil {
			return "", err
		}
		if val {
			return i.Branches[j].Run(c)
		}
	}
	if len(i.Branches) > len(i.Conditions) {
		return i.Branches[len(i.Conditions)].Run(c)
	}
	return "", nil
}

// TryBlock is a try-catch block.
type TryBlock struct {
	Try      Block
	Catch    Block
	Variable *Argument
}

// Run executes the try-catch block.
func (t *TryBlock) Run(c Context) (string, error) {
	_, thrown := t.Try.Run(c)
	if thrown == nil {
		return "", nil
	}

	if t.Variable != nil {
		name, err := t.Variable.Run(c)
		if err != nil {
			return "", nil
		}
		if _, err := c.Run("set", []string{name, thrown.Error()}); err != nil {
			return "", err
		}
	}

	return t.Catch.Run(c)
}

// WhileBlock is a "while" loop.
type WhileBlock struct {
	Condition Condition
	Block     Block
}

// Run executes the "while" loop.
func (w *WhileBlock) Run(c Context) (string, error) {
	for {
		if val, err := w.Condition.Evaluate(c); err != nil {
			return "", err
		} else if !val {
			break
		}
		_, err := w.Block.Run(c)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}
