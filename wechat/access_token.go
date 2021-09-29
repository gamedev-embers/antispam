package wechat

// AccessToken ...
type AccessToken interface {
	GetToken(force ...bool) string
	RefreshIf(bool) bool
}
