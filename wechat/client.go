package wechat

import (
	"math/rand"
	"time"
)

var _ AccessToken = &Client{}

// Client ...
type Client struct {
	appId  string
	secret string
	rand   *rand.Rand
}

func NewClient(appId, secret string) *Client {
	return &Client{
		appId:  appId,
		secret: secret,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) GetToken(force ...bool) string {
	return ""
}

func (c *Client) RefreshIf(force bool) bool {
	return false
}
