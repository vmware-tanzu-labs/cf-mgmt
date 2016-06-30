package ldap

//Manager -
type Manager interface {
	GetUserIDs(groupName string) (users []User, err error)
	IsEnabled() bool
}

//DefaultManager -
type DefaultManager struct {
	Config Config
}

//Config -
type Config struct {
	Enabled           bool   `yaml:"enabled"`
	LdapHost          string `yaml:"ldapHost"`
	LdapPort          int    `yaml:"ldapPort"`
	BindDN            string `yaml:"bindDN"`
	BindPassword      string `yaml:"bindPwd"`
	UserSearchBase    string `yaml:"userSearchBase"`
	UserNameAttribute string `yaml:"userNameAttribute"`
	UserMailAttribute string `yaml:"userMailAttribute"`
	GroupSearchBase   string `yaml:"groupSearchBase"`
	GroupAttribute    string `yaml:"groupAttribute"`
}

//User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
