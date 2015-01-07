package pragmash

const (
	BlockTypeCommand = iota
	BlockTypeFor     = iota
	BlockTypeIf      = iota
	BlockTypeWhile   = iota
)

// ParseScript parses a script and returns an array of blocks.
func ParseScript(script string) ([]Block, error) {
	// TODO: this
	return nil, nil
}

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

// ForBlock is a "for" loop.
type ForBlock struct {
	HasVariable bool
	Variable    string
	Blocks      []Block
}

// IfBlock is an "if" statement.
type IfBlock struct {
	// An array of conditions.
	Conditions []*Condition

	// An array of code block arrays.
	// Except for the last array, which may be the "else" section, each element
	// in this array should correspond to a condition.
	// len(Blocks) must be either len(Conditions) or len(Conditions)+1.
	Blocks [][]Block
}

// Type returns BlockTypeIf
func (i *IfBlock) Type() int {
	return BlockTypeIf
}

// WhileBlock is a "while" loop.
type WhileBlock struct {
	Condition *Condition
	Blocks    []Block
}

// Type returns BlockTypeWhile
func (w *WhileBlock) Type() int {
	return BlockTypeWhile
}
