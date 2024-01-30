package request

import "fmt"

func (c *Client) GetCardAll() error {
	fmt.Println("TEST")
	resp, err := c.client.R().SetCookie(c.jwt).SetHeader("Content-Type", "application/json").
		Post("http://localhost:8080/test")
	fmt.Println(resp, err)
	fmt.Println("TEST2")
	return err
}
