package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zelas91/goph-keeper/internal/client/commands"
	"github.com/zelas91/goph-keeper/internal/client/request"
	"github.com/zelas91/goph-keeper/internal/client/session"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Client struct {
	cm         *commands.CommandManager
	in         *bufio.Reader
	work       chan os.Signal
	session    *session.Session
	httpClient *resty.Client
	auth       *request.Authorization
}

func NewClient(addr string) *Client {
	cl := resty.New()
	cl.SetTimeout(time.Second * 5)
	return &Client{
		cm:         commands.NewCommandManager(),
		in:         bufio.NewReader(os.Stdin),
		work:       make(chan os.Signal),
		session:    session.NewSession(addr),
		httpClient: cl,
	}
}

func (c *Client) Start() {
	c.init()
	signal.Notify(c.work, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for {
			if !c.session.IsAuth() {
				fmt.Print("(no authorization)")
			}
			fmt.Println(">")
			args, err := commandParsing(c.in)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = c.cm.ExecCommandWithName(args[0], args[1:])
			if errors.Is(err, os.ErrProcessDone) {
				c.work <- syscall.SIGQUIT
				break
			}
			if err != nil {
				fmt.Println("command err ", err)
			}

		}
	}()
	<-c.work
}
func (c *Client) init() {
	c.auth = request.NewAuthorization(c.httpClient, c.session)
	c.registerCommandAuth()
}

func (c *Client) registerCommandAuth() {
	c.cm.RegisterCommand("login", "login into the program with an existing user",
		c.auth.SignIn, "login <login> <password>")
	c.cm.RegisterCommand("registration", "new user and login",
		c.auth.SignUp, "login <login> <password>")
}
func commandParsing(in *bufio.Reader) ([]string, error) {
	choice, err := in.ReadString('\n')
	if err != nil {
		return nil, err
	}
	choice = strings.Trim(choice, "\r\n")
	if len(choice) < 1 {
		return nil, fmt.Errorf("empty choice")
	}

	return strings.Split(choice, " "), nil

}
