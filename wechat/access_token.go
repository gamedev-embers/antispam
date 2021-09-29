package wechat

// AccessToken ...
type AccessToken interface {
	GetToken() string
	RefreshIf() bool
}
