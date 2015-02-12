package pragmash

import (
	"errors"
)

// An If object represents an "if" statement.
type If struct {
	Branches   []Runnable
	Conditions []Runnable
}

// Run executes the if statement.
func (block If) Run(r Runner) (Value, *Exception) {
	for i, cond := range block.Conditions {
		val, exc := cond.Run(r)
		if exc != nil {
			return nil, exc
		}
		if val.Bool() {
			// Run the branch.
			return block.Branches[i].Run(r)
		}
	}
	// No branches ran.
	return StringValue(""), nil
}

// An IfScanner scans an if-statement with its accompanying "else if" and
// "else" blocks.
type IfScanner struct {
	branches    []Runnable
	conditions  []Runnable
	lastContext string
	readingElse bool
	scanner     SemanticScanner
}

// NewIfScanner creates an IfScanner or fails if the initiating line is invalid.
func NewIfScanner(l Line, context string) (*IfScanner, error) {
	if len(l.Tokens) < 1 || l.Tokens[0].String != "if" {
		return nil, errors.New("if block starts with 'if' token")
	} else if l.Close || !l.Open {
		return nil, errors.New("first if line must open a block but not " +
			"close one")
	}

	condition := ConditionFromTokens(l.Tokens[1:], context)
	return &IfScanner{[]Runnable{}, []Runnable{condition}, context, false,
		newGenericScanner(true)}, nil
}

// EOF returns an error with the context of the last branch initiator.
func (s *IfScanner) EOF() (Runnable, error) {
	return nil, errors.New("missing '}' for branch at " + s.lastContext)
}

// Line adds a line to the if statement.
// If the line terminates the statement, this returns it as Runnable.
// If any kind of error is encountered, this returns the error.
// If the statement is not terminated and the line is properly processed, this
// returns nil, nil.
func (s *IfScanner) Line(l Line, context string) (Runnable, error) {
	res, err := s.scanner.Line(l, context)
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, nil
	}

	s.branches = append(s.branches, res)

	// If this closed an else block, we are definitely done.
	if s.readingElse {
		if len(l.Tokens) != 0 || l.Open {
			return nil, errors.New("else block must not be followed by any " +
				"tokens")
		}
		return s.result(), nil
	}

	// Return if we're done (i.e. this line does not open anything)
	if len(l.Tokens) == 0 && l.Open {
		return nil, errors.New("cannot open new branch without arguments")
	} else if len(l.Tokens) != 0 && !l.Open {
		return nil, errors.New("unexpected token(s) after '}'")
	} else if !l.Open {
		return s.result(), nil
	}

	// Setup the next branch
	if l.Tokens[0].String != "else" {
		return nil, errors.New("unexpected token after '}'")
	}
	s.scanner = newGenericScanner(true)
	s.lastContext = context
	if len(l.Tokens) == 1 {
		s.readingElse = true
		cond := ValueRunnable{BoolValue(true)}
		s.conditions = append(s.conditions, cond)
		return nil, nil
	} else if l.Tokens[1].String != "if" {
		return nil, errors.New("expected 'if' following 'else'")
	} else {
		cond := ConditionFromTokens(l.Tokens[2:], context)
		s.conditions = append(s.conditions, cond)
		return nil, nil
	}
}

func (s *IfScanner) result() Runnable {
	return If{s.branches, s.conditions}
}
