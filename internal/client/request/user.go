package request

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"net/http"
	"strings"
)

func (c *Client) SignIn(user models.User) error {
	resp, err := c.client.R().SetBody(user).Post(fmt.Sprintf("%s/api/signin", c.url))
	if err != nil {
		fmt.Println("err ", err)
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("ошибка запроса %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			c.jwt = cookie
			break
		}
	}
	fmt.Println("Вход выполнин")
	return nil
}

func (c *Client) SignUp(user models.User) error {
	resp, err := c.client.R().SetBody(user).Post(fmt.Sprintf("%s/api/signup", c.url))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("ошибка запроса %s", string(resp.Body()))
	}
	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, "jwt") {
			c.jwt = cookie
			break
		}
	}
	fmt.Println("Регистрация успешна, вы вошли")
	return nil
}
