package client

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/zelas91/goph-keeper/internal/client/commands"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/client/request"
	"github.com/zelas91/goph-keeper/internal/client/session"
)

type Client struct {
	cm         *commands.CommandManager
	in         *bufio.Reader
	work       chan os.Signal
	session    *session.Session
	httpClient *resty.Client
	auth       *request.Authorization
	binary     *request.BinaryFile
	card       *request.Card
	text       *request.TextData
	credential *request.Credential
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
			fmt.Print(">")
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
				if errors.Is(err, error2.ErrAuthorization) {
					c.session.CleanToken()
				}
				fmt.Println("command err ", err)
				continue
			}
			fmt.Println("success")
		}
	}()
	<-c.work
}
func (c *Client) init() {
	r := request.NewRequest(c.httpClient, c.session)

	c.auth = request.NewAuthorization(r)
	c.binary = request.NewBinaryFile(r)
	c.card = request.NewCard(r)
	c.text = request.NewTextData(r)
	c.credential = request.NewCredential(r)

	c.registerCommandAuth()
	c.registerCommandBinaryFile()
	c.registerCommandCard()
	c.registerCommandText()
	c.registerCommandCredential()
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
	c.cm.RegisterCommand("file_delete", "delete file from server",
		c.binary.Delete, "file_delete: <id>", tag)
	c.cm.RegisterCommand("files", "get data about files  on the server",
		c.binary.Files, "", tag)

	c.cm.RegisterCommand("file_download", "download file from server",
		c.binary.Download, "file_download: <id> <path>", tag)

	c.cm.RegisterCommand("file_upload", "upload file in server",
		c.binary.Upload, "file_upload: <name> <path>", tag)
}

func (c *Client) registerCommandCard() {
	tag := "credit card"
	c.cm.RegisterCommand("card_delete", "delete card from server",
		c.card.Delete, "card_delete: <id>", tag)
	c.cm.RegisterCommand("cards", "get data about cards  on the server",
		c.card.Cards, "", tag)

	c.cm.RegisterCommand("card_create", "create card to server",
		c.card.Create, "card_create: <number> <expired> <cvv>", tag)

	c.cm.RegisterCommand("card_update", "update card to server",
		c.card.Update, "card_update: <id> and any fields in the format <number:1234> <expired:12/27> <cvv:567>", tag)
}

func (c *Client) registerCommandText() {
	tag := "Text"
	c.cm.RegisterCommand("text_delete", "delete text from server",
		c.text.Delete, "text_delete: <id>", tag)
	c.cm.RegisterCommand("texts", "get data about texts  on the server",
		c.text.Texts, "", tag)

	c.cm.RegisterCommand("text_create", "create text to server",
		c.text.Create, "text_create: <'text'>", tag)

	c.cm.RegisterCommand("text_update", "update text to server",
		c.text.Update, "text_update: <id> <'text'>", tag)
}

func (c *Client) registerCommandCredential() {
	tag := "User credential"
	c.cm.RegisterCommand("credential_delete", "delete credential from server",
		c.credential.Delete, "credential_delete: <id>", tag)
	c.cm.RegisterCommand("credentials", "get data about credentials  on the server",
		c.credential.Credentials, "", tag)

	c.cm.RegisterCommand("credential_create", "create credential to server",
		c.credential.Create, "credential_create: <login> <password>", tag)

	c.cm.RegisterCommand("credential_update", "update credential to server",
		c.credential.Update, "credential_update: <id> and any fields in the format <login:password>  <password:password>", tag)
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
