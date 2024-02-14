package request

import (
	"fmt"
	"net/http"
	"strings"

	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/server/models"
)

type Authorization struct {
	request *Request
}

func NewAuthorization(request *Request) *Authorization {
	return &Authorization{request: request}
}

func (a *Authorization) SignIn(args []string) error {
	user, err := newUserModels(args)
	if err != nil {
		return err
	}

	resp, err := a.request.R().SetBody(user).Post("/signin")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response fault %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			a.request.SetCookiesAuthorization(cookie)
			break
		}
	}
	return nil
}
func (a *Authorization) SignUp(args []string) error {
	user, err := newUserModels(args)
	if err != nil {
		return err
	}
	resp, err := a.request.R().SetBody(user).Post("/signup")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response fault %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			a.request.SetCookiesAuthorization(cookie)
			break
		}
	}
	return nil
}
func newUserModels(args []string) (*models.User, error) {
	if len(args) < 2 {
		return nil, error2.ErrInvalidCommand
	}
	return &models.User{Login: args[0], Password: args[1]}, nil
}
