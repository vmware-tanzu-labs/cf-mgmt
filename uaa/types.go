package uaa

//Manager -
type Manager interface {
	GetToken() (token string, err error)
}

//Token -
type Token struct {
	AccessToken string `json:"access_token"`
}

//DefaultUAAManager -
type DefaultUAAManager struct {
	SysDomain string
	UserID    string
	Password  string
}
