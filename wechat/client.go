package wechat

import (
	"math/rand"
	"time"
)

var _ AccessToken = Client{}

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

func (this *Client) GetToken() string {
	return ""
}

func (this *Client) RefreshIf() bool {
	return false
}
