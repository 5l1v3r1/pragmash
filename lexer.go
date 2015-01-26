package pragmash

import (
	"errors"
	"strconv"
)

func ParseProgram(scriptStr string) (Blocks, error) {
	script, err := ParseScript(scriptStr)
	if err != nil {
		return nil, err
	}

	// Let's just make the line numbers start at 1.
	for i, x := range script.LineStarts {
		script.LineStarts[i] = x + 1
	}

	// Read all the blocks in the script and return it.
	res := Blocks{}
	ctx := parseContext{script, 0}
	for !ctx.done() {
		next, err := ctx.nextBlock()
		if err != nil {
			return nil, err
		} else if next == nil {
			continue
		}
		res = append(res, next)
	}
	return res, nil
}

type parseContext struct {
	script  *Script
	current int
}

func (p *parseContext) done() bool {
	return p.current >= len(p.script.LogicalLines)
}

func (p *parseContext) nextBlock() (Block, error) {
	if p.current == len(p.script.LogicalLines) {
		return nil, nil
	}

	// This prefix will be used from all errors which are returned from this
	// block.
	errorPrefix := "From line " +
		strconv.Itoa(p.script.LineStarts[p.current]) + ": "

	// Tokenize the next line
	tokens, err := Tokenize(p.script.LogicalLines[p.current])
	if err != nil {
		return nil, errors.New(errorPrefix + err.Error())
	}
	p.current++

	if len(tokens) == 0 {
		return nil, nil
	}

	// Handle control blocks.
	if tokens[0].Command == nil {
		name := tokens[0].Text
		var special Block
		var err error
		if name == "for" {
			special, err = p.readForLoop(tokens)
		} else if name == "if" {
			special, err = p.readIf(tokens)
		} else if name == "try" {
			special, err = p.readTryCatch(tokens)
		} else if name == "while" {
			special, err = p.readWhileLoop(tokens)
		}
		if err != nil {
			return nil, errors.New(errorPrefix + err.Error())
		} else if special != nil {
			return special, nil
		}
	}

	// Read a regular command and return it.
	if cmd, err := tokensToCommand(tokens); err != nil {
		return nil, errors.New(errorPrefix + err.Error())
	} else {
		return cmd, nil
	}
}

func (p *parseContext) readBlockBody(allowExtra bool) (Blocks, error) {
	res := Blocks{}
	for !p.done() {
		// Attempt to parse the line to check if it's a close curly-brace.
		tokens, err := Tokenize(p.script.LogicalLines[p.current])
		if err == nil && len(tokens) > 0 {
			if tokens[0].Text == "}" {
				if !allowExtra && len(tokens) > 1 {
					lineNum := p.script.LineStarts[p.current]
					return nil, errors.New("Unexpected tokens after } " +
						"on line " + strconv.Itoa(lineNum))
				}
				p.current++
				return res, nil
			}
		}

		// Read the next line as a block
		next, err := p.nextBlock()
		if err != nil {
			return nil, err
		} else if next == nil {
			continue
		}
		res = append(res, next)
	}
	return nil, errors.New("Missing } (at EOF).")
}

func (p *parseContext) readForLoop(t []Token) (Block, error) {
	if !endsWithOpenCurly(t) {
		return nil, errors.New("Missing { in for-loop.")
	} else if len(t) != 3 && len(t) != 4 {
		return nil, errors.New("Invalid number of arguments for for-loop.")
	}

	// Parse the arguments to the loop.
	args := make([]Argument, len(t)-2)
	for i := 1; i < len(t)-1; i++ {
		arg, err := tokenToArgument(t[i])
		if err != nil {
			return nil, err
		}
		args[i-1] = *arg
	}

	// Read the body of the loop.
	body, err := p.readBlockBody(false)
	if err != nil {
		return nil, err
	}

	// Return the for block.
	if len(args) == 1 {
		return &ForBlock{nil, args[0], body}, nil
	} else {
		return &ForBlock{&args[0], args[1], body}, nil
	}
}

func (p *parseContext) readIf(t []Token) (Block, error) {
	if !endsWithOpenCurly(t) {
		return nil, errors.New("Missing { in if-statement.")
	} else if len(t) == 2 {
		return nil, errors.New("Missing conditional in if-statement.")
	}

	res := IfBlock{make([]Condition, 1), make([]Block, 1)}

	// Read the first condition and body
	var err error
	res.Conditions[0], err = tokensToCondition(t[1 : len(t)-1])
	if err != nil {
		return nil, err
	}
	res.Branches[0], err = p.readBlockBody(true)
	if err != nil {
		return nil, err
	}

	// Read the rest of the "else if" and "else" sections.
	for {
		tokens, _ := Tokenize(p.script.LogicalLines[p.current-1])
		errLine := strconv.Itoa(p.script.LineStarts[p.current-1])

		// Handle trivial errors and endings.
		if len(tokens) == 1 {
			break
		} else if !endsWithOpenCurly(tokens) || len(tokens) == 2 ||
			len(tokens) == 4 {
			return nil, errors.New("Unexpected tokens after } on line " +
				errLine)
		}

		// Read else clause.
		if len(tokens) == 3 {
			if tokens[1].Text != "else" {
				return nil, errors.New("Unexpected token after } on line " +
					errLine)
			}
			elseBlock, err := p.readBlockBody(false)
			if err != nil {
				return nil, err
			}
			res.Branches = append(res.Branches, elseBlock)
			break
		}

		// Make sure it's an "else if"
		if tokens[1].Text != "else" || tokens[2].Text != "if" {
			return nil, errors.New("Unexpected tokens after } on line " +
				errLine)
		}
		condition, err := tokensToCondition(tokens[3 : len(tokens)-1])
		if err != nil {
			return nil, errors.New("Invalid condition on line " + errLine +
				": " + err.Error())
		}
		branch, err := p.readBlockBody(true)
		if err != nil {
			return nil, err
		}
		res.Conditions = append(res.Conditions, condition)
		res.Branches = append(res.Branches, branch)
	}

	return &res, nil
}

func (p *parseContext) readTryCatch(t []Token) (Block, error) {
	if !endsWithOpenCurly(t) {
		return nil, errors.New("Missing { in try-catch block.")
	} else if len(t) != 2 {
		return nil, errors.New("Invalid extra arguments for try-catch block.")
	}

	// Read the body of the try block
	body, err := p.readBlockBody(true)
	if err != nil {
		return nil, err
	}

	// Handle empty or invalid catch blocks.
	lastLine, _ := Tokenize(p.script.LogicalLines[p.current-1])
	invalMessage := errors.New("Invalid arguments after } on line " +
		strconv.Itoa(p.script.LineStarts[p.current-1]))
	if len(lastLine) == 1 {
		// No catch block.
		return &TryBlock{body, Blocks{}, nil}, nil
	} else if len(lastLine) != 3 && len(lastLine) != 4 {
		return nil, invalMessage
	}

	// Line must be "} catch [variable] {"
	if !endsWithOpenCurly(lastLine) || lastLine[1].Text != "catch" {
		return nil, invalMessage
	}

	// Get the optional variable argument.
	var variable *Argument
	if len(lastLine) == 4 {
		variable, err = tokenToArgument(lastLine[2])
		if err != nil {
			return nil, errors.New("Invalid variable argument at line " +
				strconv.Itoa(p.script.LineStarts[p.current-1]) + ": " +
				err.Error())
		}
	}

	// Read the catch body
	catchBody, err := p.readBlockBody(false)
	if err != nil {
		return nil, err
	}

	return &TryBlock{body, catchBody, variable}, nil
}

func (p *parseContext) readWhileLoop(t []Token) (Block, error) {
	if !endsWithOpenCurly(t) {
		return nil, errors.New("Missing { in while-loop.")
	}

	// Parse the condition.
	cond, err := tokensToCondition(t[1 : len(t)-1])
	if err != nil {
		return nil, err
	}

	// Read the body of the loop.
	body, err := p.readBlockBody(false)
	if err != nil {
		return nil, err
	}

	return &WhileBlock{cond, body}, nil
}

func endsWithOpenCurly(t []Token) bool {
	if len(t) == 1 {
		return false
	}
	last := t[len(t)-1]
	return last.Text == "{"
}

func tokenToArgument(t Token) (*Argument, error) {
	if t.Command == nil {
		return &Argument{t.Text, nil}, nil
	}

	// Parse the sub-command.
	command, err := tokensToCommand(t.Command)
	if err != nil {
		return nil, err
	}
	return &Argument{"", command}, nil
}

func tokensToCommand(t []Token) (*Command, error) {
	if len(t) == 0 {
		return nil, errors.New("No tokens in command.")
	}
	args := make([]Argument, len(t))
	for i, x := range t {
		arg, err := tokenToArgument(x)
		if err != nil {
			return nil, err
		}
		args[i] = *arg
	}
	return &Command{args[0], args[1:]}, nil
}

func tokensToCondition(t []Token) (Condition, error) {
	args := make(Condition, len(t))
	for i, x := range t {
		arg, err := tokenToArgument(x)
		if err != nil {
			return nil, err
		}
		args[i] = *arg
	}
	return args, nil
}
