package uaac

//Manager -
type Manager interface {

	//Returns a map keyed and valued by user id. User id is converted to lowercase
	ListUsers() (map[string]string, error)

	// Returns a map keyed by userid and value as User struct.
	// Return an empty map if an error occurs or if there are no users found
	UsersByID() (map[string]User, error)

	CreateExternalUser(userName, userEmail, externalID, origin string) (err error)
}

//UserList -
type UserList struct {
	Users []User `json:"resources"`
}

//User -
type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Origin   string `json:"origin"`
}
