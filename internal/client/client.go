package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zelas91/goph-keeper/internal/client/commands"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
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
	binary     *request.BinaryFile
}

func NewClient(addr string) *Client {
	cl := resty.New()
	cl.SetTimeout(time.Second * 5)
	cl.SetBaseURL(fmt.Sprintf("http://%s/api", addr))
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
				if errors.Is(err, error2.ErrAuthorization) {
					c.session.CleanToken()
				}
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
	r := request.NewRequest(c.httpClient, c.session)
	c.auth = request.NewAuthorization(r)
	c.binary = request.NewBinaryFile(r)
	c.registerCommandBinaryFile()
	c.registerCommandAuth()
}

func (c *Client) registerCommandAuth() {
	tag := "authorization"
	c.cm.RegisterCommand("login", "login into the program with an existing user",
		c.auth.SignIn, "login: <login> <password>", tag)
	c.cm.RegisterCommand("registration", "new user and login",
		c.auth.SignUp, "login: <login> <password>", tag)
}

func (c *Client) registerCommandBinaryFile() {
	tag := "Binary File"
	c.cm.RegisterCommand("delete_file", "delete file from server",
		c.binary.DeleteFile, "delete_file: <id>", tag)
	c.cm.RegisterCommand("get_files", "get data about files  on the server",
		c.binary.GetAll, "", tag)

	c.cm.RegisterCommand("download_file", "download file from server",
		c.binary.Download, "download_file: <id> <path>", tag)

	c.cm.RegisterCommand("upload_file", "upload file in server",
		c.binary.Upload, "upload_file: <name> <path>", tag)
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
