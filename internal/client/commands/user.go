package commands

import (
	"bufio"
	"fmt"
	"github.com/zelas91/goph-keeper/internal/client"
	"github.com/zelas91/goph-keeper/internal/server/models"
)

type authorization struct {
	in          *bufio.Reader
	authService userService
}

type userService interface {
	SignIn(user models.User) error
	SignUp(user models.User) error
	GetCardAll() error
}

func (u *authorization) registerAndLoginCommand() error {
	fmt.Println("Выберите: 1 - Регистрация, 2 - Вход")
	command, err := commandInt(u.in)

	if err != nil {
		return &client.ErrAuth{Err: fmt.Errorf("не вверный ввод %v", err)}
	}

	switch command {
	case 1:
		if err = u.signUp(); err != nil {
			return &client.ErrAuth{Err: err}
		}
	case 2:
		if err = u.signIn(); err != nil {
			return &client.ErrAuth{Err: err}
		}
	default:
		return &client.ErrAuth{Err: fmt.Errorf("выбрана не верная команда %d", command)}
	}
	return nil
}

func (u *authorization) signIn() error {
	user, err := newUserModels(u.in)
	if err != nil {
		return err
	}
	return u.authService.SignIn(user)
}

func (u *authorization) signUp() error {
	user, err := newUserModels(u.in)
	if err != nil {
		return err
	}
	return u.authService.SignUp(user)

}

func newUserModels(in *bufio.Reader) (models.User, error) {
	var user models.User
	var err error
	fmt.Println("введите логин:")
	user.Login, err = commandStr(in)
	if err != nil {
		fmt.Printf("ошибка ввода логина %v\n", err)
		return user, err
	}
	fmt.Println("введите пароль:")
	user.Password, err = commandStr(in)
	if err != nil {
		fmt.Printf("ошибка ввода пароля %v\n", err)
		return user, err
	}
	return user, err
}
