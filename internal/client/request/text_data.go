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

type TextData struct {
	request *Request
}

func NewTextData(request *Request) *TextData {
	return &TextData{request: request}
}
func (t *TextData) Delete(args []string) error {
	if len(args) < 1 {
		return error2.ErrInvalidCommand
	}
	url := fmt.Sprintf("/text/%s", args[0])
	resp, err := t.request.R().Delete(url)
	if err != nil {
		return fmt.Errorf("request text delete err: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request text delete error status code = %d", resp.StatusCode())
	}
	return nil
}

func (t *TextData) Update(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}
	textID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("text id err:%v", err)
	}

	url := fmt.Sprintf("/text/%d", textID)

	resp, err := t.request.R().Get(url)
	if err != nil {
		return fmt.Errorf("request get text err: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request get text error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	var text models.TextData
	if err := json.Unmarshal(resp.Body(), &text); err != nil {
		return fmt.Errorf("request text decode err: %v", err)
	}
	var tb strings.Builder
	for i := 1; i < len(args); i++ {
		tb.WriteString(args[i])
		tb.WriteString(" ")
	}
	text.Text = tb.String()

	url = fmt.Sprintf("/text/%d", textID)

	resp, err = t.request.R().SetBody(text).Put(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request update text error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (t *TextData) Create(args []string) error {
	if len(args) < 1 {
		return error2.ErrInvalidCommand
	}
	var tb strings.Builder
	for i := 0; i < len(args); i++ {
		tb.WriteString(args[i])
		tb.WriteString(" ")
	}
	text := models.TextData{
		Text: tb.String(),
	}

	resp, err := t.request.R().SetBody(text).Post("/text")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("request create text error status code = %d, body = %s",
			resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

func (t *TextData) Texts(args []string) error {
	resp, err := t.request.R().Get("/text")
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
