package commands

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	error2 "github.com/zelas91/goph-keeper/internal/client/error"
)

type command struct {
	Name        string
	Description string
	Exec        func(args []string) error
	help        string
	tag         string
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
	tag string,

) {
	cmd := command{
		Name:        name,
		Description: description,
		Exec:        exec,
		help:        help,
		tag:         tag,
	}

	c.commands[name] = cmd
}

func (c *CommandManager) ExecCommandWithName(name string, args []string) error {
	cmd, ok := c.commands[name]
	if !ok {
		return errors.New("invalid command, 'help'")
	}

	err := cmd.Exec(args)

	if errors.Is(err, error2.ErrInvalidCommand) {
		fmt.Println("help command:", cmd.help)
		return nil
	}

	return err
}

func (c *CommandManager) initBaseCommands() {
	tag := "system"
	c.RegisterCommand("help", "list all commands", c.help, "help", tag)
	c.RegisterCommand("exit", "close program", c.exit, "exit", tag)
}

func (c *CommandManager) help(args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "-----------------------------------------------------------")

	com := make(map[string][]string)

	for _, v := range c.commands {
		com[v.tag] = append(com[v.tag],
			fmt.Sprintf("command: '%s'\t Description: %s.\t use: %s",
				v.Name, v.Description, v.help))
	}
	for n, v := range com {
		fmt.Fprintln(w, "type", n)

		for _, command := range v {
			fmt.Fprintln(w, "\t", command)
		}
		fmt.Fprintln(w, "===========================================================")
	}
	fmt.Fprintln(w, "-----------------------------------------------------------")
	if err := w.Flush(); err != nil {
		return fmt.Errorf("table format err :%v", err)
	}
	return nil
}

func (c *CommandManager) exit(args []string) error {
	return os.ErrProcessDone
}
