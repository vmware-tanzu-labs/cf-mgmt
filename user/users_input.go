package user

import (
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
)

//UsersInput
type UsersInput struct {
	SpaceGUID                                   string
	OrgGUID                                     string
	LdapUsers, Users, LdapGroupNames, SamlUsers []string
	SpaceName                                   string
	OrgName                                     string
	RemoveUsers                                 bool
	OrgUsers                                    *RoleUsers
	ListUsers                                   func(updateUserInput UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error)
	AddUser                                     func(updateUserInput UsersInput, userName, userGUID string) error
	RemoveUser                                  func(updateUserInput UsersInput, userName, userGUID string) error
}

func (u *UsersInput) UniqueUsers() []string {
	return uniqueSlice(u.Users)
}

func (u *UsersInput) UniqueSamlUsers() []string {
	return uniqueSlice(u.SamlUsers)
}

func (u *UsersInput) UniqueLdapUsers() []string {
	return uniqueSlice(u.LdapUsers)
}

func (u *UsersInput) UniqueLdapGroupNames() []string {
	return uniqueSlice(u.LdapGroupNames)
}

func uniqueSlice(input []string) []string {
	unique := make(map[string]string)
	output := []string{}
	for _, value := range input {
		if _, ok := unique[value]; !ok {
			unique[value] = value
		}
	}
	for key := range unique {
		output = append(output, key)
	}
	return output
}
