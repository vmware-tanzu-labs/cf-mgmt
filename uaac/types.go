package uaac

//Manager -
type Manager interface {
	ListUsers() (users map[string]string, err error)
	CreateLdapUser(userName, userEmail, userDN string) (err error)
}

//UserList -
type UserList struct {
	Users []User `json:"resources"`
}

//User -
type User struct {
	ID   string `json:"id"`
	Name string `json:"userName"`
}
