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

type Card struct {
	request *Request
}

func NewCard(request *Request) *Card {
	return &Card{request: request}
}
func (c *Card) Delete(args []string) error {
	if len(args) < 1 {
		return error2.ErrInvalidCommand
	}
	url := fmt.Sprintf("/card/%s", args[0])
	resp, err := c.request.R().Delete(url)
	if err != nil {
		return fmt.Errorf("request card delete err: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request card delete error status code = %d", resp.StatusCode())
	}
	return nil
}

func (c *Card) Update(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}
	cardID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("card id err:%v", err)
	}

	var card models.Card
	var number, cvv, expired string

	for i := 1; i < len(args); i++ {
		values := strings.Split(args[i], ":")
		if len(values) != 2 {
			return error2.ErrInvalidCommand
		}
		switch values[0] {
		case "number":
			number = values[1]
		case "cvv":
			cvv = values[1]
		case "expired":
			expired = values[1]
		}
	}
	url := fmt.Sprintf("/card/%d", cardID)
	resp, err := c.request.R().Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request get card id error status code = %d, body=%v",
			resp.StatusCode(), string(resp.Body()))
	}

	if err := json.Unmarshal(resp.Body(), &card); err != nil {
		return fmt.Errorf("request card decode err: %v", err)
	}

	card = updateModelCard(card, number, cvv, expired)
	resp, err = c.request.R().SetBody(card).Put(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request update card error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func updateModelCard(card models.Card, number, cvv, expired string) models.Card {
	if number != "" {
		card.Number = number
	}

	if cvv != "" {
		card.Cvv = cvv
	}
	if expired != "" {
		card.ExpiredAt = expired
	}
	return card
}
func (c *Card) Create(args []string) error {
	if len(args) < 3 {
		return error2.ErrInvalidCommand
	}

	card := models.Card{
		Number:    args[0],
		ExpiredAt: args[1],
		Cvv:       args[2],
	}
	resp, err := c.request.R().SetBody(card).Post("/card")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("request create card error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (c *Card) Cards(args []string) error {
	resp, err := c.request.R().Get("/card")
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
