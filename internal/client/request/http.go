package request

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/client/session"
	"net/http"
	"net/url"
)

type Request struct {
	httClient *resty.Client
	session   *session.Session
}
type Query struct {
	req *resty.Request
}

func NewRequest(httClient *resty.Client, session *session.Session) *Request {
	return &Request{
		httClient: httClient,
		session:   session,
	}
}
func (q *Query) SetBody(body interface{}) *Query {
	q.req.SetBody(body)
	return q
}
func (q *Query) Post(url string) (*resty.Response, error) {
	return q.isAuthorization(q.req.Post(url))
}

func (q *Query) Get(url string) (*resty.Response, error) {
	return q.isAuthorization(q.req.Get(url))
}

func (q *Query) Delete(url string) (*resty.Response, error) {
	return q.isAuthorization(q.req.Delete(url))
}

func (q *Query) Put(url string) (*resty.Response, error) {
	return q.isAuthorization(q.req.Put(url))
}

func (q *Query) isAuthorization(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return resp, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return resp, error2.ErrAuthorization
	}
	return resp, err
}

func (r *Request) R() *Query {
	return &Query{
		req: r.httClient.R().
			SetCookie(r.session.GetJwt()).
			SetHeader("Content-type", "application/json"),
	}
}

func (r *Request) WebsocketConnect(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws",
		Host: r.session.Host,
		Path: fmt.Sprintf("/api%s", addr)}
	cfg := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("jwt", r.session.GetJwt().Value)
	conn, _, err := cfg.Dial(u.String(), header)
	return conn, err
}

func (r *Request) SetCookiesAuthorization(aut *http.Cookie) {
	r.session.Jwt = aut
}
