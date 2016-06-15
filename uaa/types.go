package uaa

//Manager -
type Manager interface {
	GetCFToken(password string) (token string)
	GetUAACToken(secret string) (token string)
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	SysDomain string
	UserID    string
}
