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
