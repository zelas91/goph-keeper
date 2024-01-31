package request

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/client/session"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"net/http"
	"strings"
)

type Authorization struct {
	httpClient *resty.Client
	session    *session.Session
}

func NewAuthorization(httpClient *resty.Client, session *session.Session) *Authorization {
	return &Authorization{httpClient: httpClient, session: session}
}

func (a *Authorization) SignIn(args []string) error {
	user, err := newUserModels(args)
	if err != nil {
		return err
	}
	resp, err := a.httpClient.R().SetBody(user).Post(fmt.Sprintf("%s/api/signin", a.session.Url))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response fault %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			a.session.Jwt = cookie
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
	resp, err := a.httpClient.R().SetBody(user).Post(fmt.Sprintf("%s/api/signup", a.session.Url))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response fault %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			a.session.Jwt = cookie
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
