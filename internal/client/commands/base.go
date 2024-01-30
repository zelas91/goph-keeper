package commands

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/zelas91/goph-keeper/internal/client"
	"github.com/zelas91/goph-keeper/internal/logger"
	"os"
	"strconv"
	"strings"
)

type Command struct {
	*user
	in *bufio.Reader
}

func (c *Command) Start() {
	var auth bool
	for {
		if !auth {
			err := c.user.registerAndLoginCommand()
			auth = err == nil
			if err != nil {
				fmt.Println(err)
			}

		}

		if auth {
			if err := c.workingWithData(); err != nil {
				fmt.Println(err)
				var target *client.ErrAuth
				if errors.As(err, target) {
					auth = false
				}
			}
		}
	}

}

func New(log logger.Logger, options ...func(c *Command)) *Command {
	ctl := &Command{in: bufio.NewReader(os.Stdin)}
	for _, opt := range options {
		opt(ctl)
	}
	return ctl
}

func (c *Command) workingWithData() error {
	fmt.Println("выбирити тип данных 1-card, 2-file, 3-login and password, 4-text")
	command, err := commandInt(c.in)
	if err != nil {
		return err
	}

	switch command {
	case 1:
		return c.s.GetCardAll()
	default:
		return errors.New("не верный выбор операции")
	}
}

func WithUserCommand(us userService) func(c *Command) {
	return func(c *Command) {
		c.user = &user{s: us, in: c.in}
	}
}

func commandInt(in *bufio.Reader) (int, error) {
	choice, err := in.ReadString('\n')
	if err != nil {
		return 0, err
	}
	choice = strings.TrimSuffix(choice, "\r\n")
	return strconv.Atoi(choice)

}

func commandStr(in *bufio.Reader) (string, error) {
	choice, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}
	choice = strings.TrimSuffix(choice, "\r\n")
	if len(choice) < 1 {
		return "", fmt.Errorf("поле пустое")
	}
	return choice, nil

}
