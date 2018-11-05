package uaa

import "strings"

type Users struct {
	userMap map[string][]User
}

func (u *Users) Add(user User) {
	if u.userMap == nil {
		u.userMap = make(map[string][]User)
	}
	key := strings.ToLower(user.Username)
	existingUsers := u.userMap[key]
	existingUsers = append(existingUsers, user)
	u.userMap[key] = existingUsers
}

func (u *Users) List() []User {
	if u.userMap == nil {
		return nil
	}
	var result []User
	for key := range u.userMap {
		result = append(result, u.userMap[key]...)
	}
	return result
}

func (u *Users) Exists(userName string) bool {
	if u.userMap == nil {
		return false
	}
	_, ok := u.userMap[strings.ToLower(userName)]
	return ok
}

func (u *Users) GetByName(userName string) []User {
	if u.userMap == nil {
		return nil
	}
	return u.userMap[strings.ToLower(userName)]
}

func (u *Users) GetByID(ID string) *User {
	for _, user := range u.List() {
		if strings.EqualFold(user.GUID, ID) {
			return &user
		}
	}
	return nil
}

func (u *Users) GetByExternalID(externalID string) *User {
	for _, user := range u.List() {
		if strings.EqualFold(user.ExternalID, externalID) {
			return &user
		}
	}
	return nil
}
