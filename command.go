package pragmash

import "errors"

type Command []Expression

func ParseCommand(line string) (Command, error) {
	tokens, err := Tokenize(line)
	if err != nil {
		return nil, err
	} else if len(tokens) == 0 {
		return nil, errors.New("Empty token array")
	}
	args := make([]Expression, len(tokens))
	for i, x := range tokens {
		if x.Command {
			args[i], err = ParseCommand(x.Text)
			if err != nil {
				return nil, err
			}
		} else {
			args[i] = StringExpression(x.Text)
		}
	}
	return Command(args), nil
}

func (c Command) Run(ctx *Context) (string, error) {
	return "", nil
}

type StringExpression string

func (s StringExpression) Run(c *Context) (string, error) {
	return string(s), nil
}
