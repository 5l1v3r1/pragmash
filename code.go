package pragmash

import (
	"errors"
	"strings"
)

const (
	BlockTypeCommand = iota
	BlockTypeIf = iota
	BlockTypeFor = iota
)

// Argument is an argument for a command.
// An argument is either a raw string or a sub-command.
type Argument struct {
	IsCommand bool
	String    string
	Command   *Command
}

// Block stores a command or a control-flow block.
type Block interface {
	Type() int
}

// Command is a block representing a single-line command.
type Command struct {
	Name      string
	Arguments []*Argument
}

// Type returns BlockTypeCommand.
func (c *Command) Type() int {
	return BlockTypeCommand
}

// Condition is a condition used for if-statements.
type Condition []*Argument

// IfBlock stores the information for an if-statement.
type IfBlock struct {
	Conditions []*Condition
	Blocks     []Block
}

// Type returns BlockTypeIf
func (i *IfBlock) Type() int {
	return BlockTypeIf
}
