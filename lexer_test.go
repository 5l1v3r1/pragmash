package pragmash

import (
	"testing"
)

func TestLexCommand(t *testing.T) {
	res, err := ParseProgram("echo yo `echo there \\`echo bro\\``")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid length of result:", len(res))
	}
	command, ok := res[0].(*Command)
	if !ok {
		t.Fatal("Block was not a *Command")
	}
	if command.Name.Command != nil || command.Name.Text != "echo" {
		t.Fatal("Invalid command name:", command.Name)
	}
	if len(command.Arguments) != 2 {
		t.Fatal("Invalid number of arguments.")
	}
	if command.Arguments[0].Command != nil ||
		command.Arguments[0].Text != "yo" {
		t.Error("Invalid first argument:", command.Arguments[0])
	}
	if command.Arguments[1].Command == nil {
		t.Fatal("Invalid second argument:", command.Arguments[1])
	}

	subCommand1 := command.Arguments[1].Command
	if subCommand1.Name.Command != nil || subCommand1.Name.Text != "echo" {
		t.Fatal("Invalid sub-command name.")
	}
	if len(subCommand1.Arguments) != 2 {
		t.Fatal("Invalid number of sub-command arguments.")
	}
	if subCommand1.Arguments[0].Command != nil ||
		subCommand1.Arguments[0].Text != "there" {
		t.Error("Invalid first sub-command argument.")
	}
	if subCommand1.Arguments[1].Command == nil {
		t.Fatal("Invalid second sub-command argument.")
	}

	subCommand2 := subCommand1.Arguments[1].Command
	if subCommand2.Name.Command != nil || subCommand2.Name.Text != "echo" {
		t.Fatal("Invalid sub-sub-command name.")
	}
	if len(subCommand2.Arguments) != 1 {
		t.Fatal("Invalid sub-sub-command arguments.")
	}
	if subCommand2.Arguments[0].Command != nil ||
		subCommand2.Arguments[0].Text != "bro" {
		t.Fatal("Invalid first sub-sub-command argument.")
	}
}

func TestLexForSimple(t *testing.T) {
	// Test failures.
	if _, err := ParseProgram("for hey {\n\n} foobar"); err == nil {
		t.Error("Trailing text after } should trigger error.")
	}
	if _, err := ParseProgram("for {\n\n}"); err == nil {
		t.Error("No argument should trigger error.")
	}
	if _, err := ParseProgram("for foo\n\n}"); err == nil {
		t.Error("Missing { should trigger error.")
	}

	res, err := ParseProgram("for `echo yo\\\\nthere` {\necho hey\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid length for result.")
	}
	block, ok := res[0].(*ForBlock)
	if !ok {
		t.Fatal("Invalid block type.")
	}
	if block.Variable != nil {
		t.Error("Invalid variable found.")
	}
	if block.Expression.Command == nil ||
		block.Expression.Command.Name.Command != nil ||
		block.Expression.Command.Name.Text != "echo" ||
		len(block.Expression.Command.Arguments) != 1 ||
		block.Expression.Command.Arguments[0].Command != nil ||
		block.Expression.Command.Arguments[0].Text != "yo\nthere" {
		t.Error("Invalid loop expression")
	}
}

func TestLexForVariable(t *testing.T) {
	// Test failures.
	if _, err := ParseProgram("for x hey {\n\n} foobar"); err == nil {
		t.Error("Trailing text after } should trigger error.")
	}
	if _, err := ParseProgram("for x foo\n\n}"); err == nil {
		t.Error("Missing { should trigger error.")
	}

	res, err := ParseProgram("for x `echo yo\\\\nthere` {\necho hey\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid length for result.")
	}
	block, ok := res[0].(*ForBlock)
	if !ok {
		t.Fatal("Invalid block type.")
	}
	if block.Variable == nil || block.Variable.Command != nil ||
		block.Variable.Text != "x" {
		t.Error("Invalid variable found.")
	}
	if block.Expression.Command == nil ||
		block.Expression.Command.Name.Command != nil ||
		block.Expression.Command.Name.Text != "echo" ||
		len(block.Expression.Command.Arguments) != 1 ||
		block.Expression.Command.Arguments[0].Command != nil ||
		block.Expression.Command.Arguments[0].Text != "yo\nthere" {
		t.Error("Invalid loop expression")
	}
}

func TestLexIf(t *testing.T) {
	// Do some invalid if statements
	if _, err := ParseProgram("if hey\necho yo"); err == nil {
		t.Error("No opening { should trigger error.")
	}
	if _, err := ParseProgram("if hey {\necho yo"); err == nil {
		t.Error("No closing } should trigger error.")
	}
	if _, err := ParseProgram("if hey {\necho yo\n} else\n"); err == nil {
		t.Error("Trailing else should trigger error.")
	}
	if _, err := ParseProgram("if {\necho yo\n}"); err == nil {
		t.Error("Empty condition should trigger error.")
	}
	if _, err := ParseProgram("if {\necho yo\n} else if {\n\n}"); err == nil {
		t.Error("Empty 'else if' condition should trigger error.")
	}
	if _, err := ParseProgram("if {\necho yo\n} else if foo\n\n}"); err == nil {
		t.Error("Missing 'else if' { should trigger error.")
	}
	if _, err := ParseProgram("if {\necho yo\n} else {\n} else {\n}"); err == nil {
		t.Error("Double else clause should trigger error.")
	}

	// Test parsing a full-blown if statement.
	res, err := ParseProgram("if hey {\necho yo\necho yay\n} else if " +
		"bob `echo bob` {\necho yoyo\n} else if " + "`testing123` {\n\n} " +
		"else {\necho yo1\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid block count.")
	}
	ifBlock, ok := res[0].(*IfBlock)
	if !ok {
		t.Fatal("Result was not if block")
	}

	if len(ifBlock.Conditions) != 3 {
		t.Fatal("Invalid number of conditions.")
	} else if len(ifBlock.Branches) != 4 {
		t.Fatal("Invalid number of branches.")
	}

	// Validate the conditions.
	if len(ifBlock.Conditions[0]) != 1 ||
		ifBlock.Conditions[0][0].Command != nil ||
		ifBlock.Conditions[0][0].Text != "hey" {
		t.Fatal("First condition was invalid.")
	} else if len(ifBlock.Conditions[1]) != 2 ||
		ifBlock.Conditions[1][0].Command != nil ||
		ifBlock.Conditions[1][0].Text != "bob" ||
		ifBlock.Conditions[1][1].Command == nil ||
		ifBlock.Conditions[1][1].Command.Name.Command != nil ||
		ifBlock.Conditions[1][1].Command.Name.Text != "echo" ||
		len(ifBlock.Conditions[1][1].Command.Arguments) != 1 ||
		ifBlock.Conditions[1][1].Command.Arguments[0].Command != nil ||
		ifBlock.Conditions[1][1].Command.Arguments[0].Text != "bob" {
		t.Fatal("Second condition was invalid.")
	} else if len(ifBlock.Conditions[2]) != 1 ||
		ifBlock.Conditions[2][0].Command == nil ||
		ifBlock.Conditions[2][0].Command.Name.Command != nil ||
		ifBlock.Conditions[2][0].Command.Name.Text != "testing123" ||
		len(ifBlock.Conditions[2][0].Command.Arguments) != 0 {
		t.Fatal("Third condition was invalid.")
	}

	// Briefly validate the block lengths.
	if len(ifBlock.Branches[0].(Blocks)) != 2 {
		t.Fatal("First branch was invalid.")
	} else if len(ifBlock.Branches[1].(Blocks)) != 1 {
		t.Fatal("Second branch was invalid.")
	} else if len(ifBlock.Branches[2].(Blocks)) != 0 {
		t.Fatal("Third branch was invalid.")
	} else if len(ifBlock.Branches[3].(Blocks)) != 1 {
		t.Fatal("Fourth branch was invalid.")
	}
}

func TestLexTryCatch(t *testing.T) {
	// Some invalid try cases.
	if _, err := ParseProgram("try {\n} catch {\n} foo"); err == nil {
		t.Error("Trailing characters after } should trigger error.")
	}
	if _, err := ParseProgram("try {\n} foo x {\n}"); err == nil {
		t.Error("Missing 'catch' should trigger error.")
	}

	res, err := ParseProgram("try {\ndie foo\n} catch {\necho $x\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid result length.")
	}

	block, ok := res[0].(*TryBlock)
	if !ok {
		t.Error("Invalid block type.")
	}

	if block.Variable != nil {
		t.Error("Expected nil variable.")
	}
	if len(block.Try.(Blocks)) != 1 {
		t.Error("Expected one try block.")
	}
	if len(block.Catch.(Blocks)) != 1 {
		t.Error("Expected one catch block.")
	}
}

func TestLexTryNoCatch(t *testing.T) {
	// Some invalid try cases.
	if _, err := ParseProgram("try {\n"); err == nil {
		t.Error("Missing } should trigger error.")
	}
	if _, err := ParseProgram("try {\n} foobar"); err == nil {
		t.Error("Trailing text after } should trigger error.")
	}
	if _, err := ParseProgram("try foo {\n}"); err == nil {
		t.Error("Extra argument should trigger error.")
	}
	if _, err := ParseProgram("try foo\n}"); err == nil {
		t.Error("Missing { should trigger error")
	}

	res, err := ParseProgram("try {\ndie foo\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid result length.")
	}

	block, ok := res[0].(*TryBlock)
	if !ok {
		t.Error("Invalid block type.")
	}

	if block.Variable != nil {
		t.Error("Expected nil variable.")
	}
	if len(block.Try.(Blocks)) != 1 {
		t.Error("Expected one try block.")
	}
	if len(block.Catch.(Blocks)) != 0 {
		t.Error("Expected empty blocks for catch block.")
	}
}

func TestLexTryVarCatch(t *testing.T) {
	// Some invalid try cases.
	if _, err := ParseProgram("try {\n} catch x y {\n}"); err == nil {
		t.Error("Multiple catch variables should trigger error.")
	}
	if _, err := ParseProgram("try {\n} catch x y\n}"); err == nil {
		t.Error("Missing { for catch block should trigger error.")
	}

	res, err := ParseProgram("try {\ndie foo\n} catch x {\necho $x\n}")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatal("Invalid result length.")
	}

	block, ok := res[0].(*TryBlock)
	if !ok {
		t.Error("Invalid block type.")
	}

	if block.Variable == nil || block.Variable.Command != nil ||
		block.Variable.Text != "x" {
		t.Error("Expected 'x' variable.")
	}
	if len(block.Try.(Blocks)) != 1 {
		t.Error("Expected one try block.")
	}
	if len(block.Catch.(Blocks)) != 1 {
		t.Error("Expected one catch block.")
	}
}
