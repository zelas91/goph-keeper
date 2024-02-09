package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/server/models"
)

type Credential struct {
	request *Request
}

func NewCredential(request *Request) *Credential {
	return &Credential{request: request}
}
func (c *Credential) Delete(args []string) error {
	if len(args) < 1 {
		return error2.ErrInvalidCommand
	}
	url := fmt.Sprintf("/credential/%s", args[0])
	resp, err := c.request.R().Delete(url)
	if err != nil {
		return fmt.Errorf("request credential delete err: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request credential delete error status code = %d", resp.StatusCode())
	}
	return nil
}

func (c *Credential) Update(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}
	credentialID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("credential id err:%v", err)
	}

	var credential models.UserCredentials
	var login, password string

	for i := 1; i < len(args); i++ {
		values := strings.Split(args[i], ":")
		if len(values) != 2 {
			return error2.ErrInvalidCommand
		}
		switch values[0] {
		case "login":
			login = values[1]
		case "password":
			password = values[1]
		}
	}
	url := fmt.Sprintf("/credential/%d", credentialID)
	resp, err := c.request.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request get credential id error status code = %d, body=%v",
			resp.StatusCode(), string(resp.Body()))
	}

	if err := json.Unmarshal(resp.Body(), &credential); err != nil {
		return fmt.Errorf("request credential decode err: %v", err)
	}

	credential = updateModelCredential(credential, login, password)
	resp, err = c.request.R().SetBody(credential).Put(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request update credential error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func updateModelCredential(credential models.UserCredentials, login, password string) models.UserCredentials {
	if login != "" {
		credential.Login = login
	}
	if password != "" {
		credential.Password = password
	}
	return credential
}
func (c *Credential) Create(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}

	credential := models.UserCredentials{
		Login:    args[0],
		Password: args[1],
	}
	resp, err := c.request.R().SetBody(credential).Post("/credential")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("request create credential error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *Credential) Credentials(args []string) error {
	resp, err := c.request.R().Get("/credential")
	if err != nil {
		return err
	}
	str, err := prettyJSON(resp.Body())
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}
