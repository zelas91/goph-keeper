package request

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const (
	baseURL = "http://"
)

type Client struct {
	client *resty.Client
	jwt    *http.Cookie
	url    string
}

func NewClient(url string) *Client {
	cl := resty.New()
	cl.SetTimeout(time.Second * 2)
	return &Client{client: cl, url: fmt.Sprintf("%s%s", baseURL, url)}
}
