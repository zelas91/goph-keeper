package commands

import (
	"errors"
	"fmt"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"os"
)

type command struct {
	Name        string
	Description string
	Exec        func(args []string) error
	help        string
}

type CommandManager struct {
	commands map[string]command
}

func NewCommandManager() *CommandManager {
	manager := CommandManager{
		commands: make(map[string]command),
	}
	manager.initBaseCommands()
	return &manager
}
func (c *CommandManager) RegisterCommand(
	name string,
	description string,
	exec func(args []string) error,
	help string,

) {
	cmd := command{
		Name:        name,
		Description: description,
		Exec:        exec,
		help:        help,
	}

	c.commands[name] = cmd
}

func (c *CommandManager) ExecCommandWithName(name string, args []string) error {
	cmd, ok := c.commands[name]
	if !ok {
		fmt.Println("invalid command, 'help'")
		return nil
	}

	err := cmd.Exec(args)

	if errors.Is(err, error2.ErrInvalidCommand) {
		fmt.Println("Usage:", cmd.help)
		return nil
	}

	return err
}

func (c *CommandManager) initBaseCommands() {
	c.RegisterCommand("help", "list all commands", c.help, "help")
	c.RegisterCommand("exit", "close program", c.exit, "exit")
}

func (c *CommandManager) help(args []string) error {
	for n, c := range c.commands {
		fmt.Println("command name:", n, "description:", c.Description)
	}
	fmt.Println("-----------------------------------------------------------")
	return nil
}
func (c *CommandManager) exit(args []string) error {
	return os.ErrProcessDone
}
