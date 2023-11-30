package user

import (
	"fmt"
	"strings"

	"github.com/vmwarepivotallabs/cf-mgmt/role"
)

// UsersInput
type UsersInput struct {
	SpaceGUID                                             string
	OrgGUID                                               string
	LdapUsers, Users, SPNUsers, LdapGroupNames, SamlUsers []string
	SpaceName                                             string
	OrgName                                               string
	RemoveUsers                                           bool
	RoleUsers                                             *role.RoleUsers
	AddUser                                               func(orgGUID, entityName, entityGUID, userName, userGUID string) error
	RemoveUser                                            func(entityName, entityGUID, userName, userGUID string) error
	Role                                                  string
}

func (u *UsersInput) EntityName() string {
	if u.SpaceGUID != "" {
		return fmt.Sprintf("%s/%s", u.OrgName, u.SpaceName)
	}
	return u.OrgName
}

func (u *UsersInput) EntityGUID() string {
	if u.SpaceGUID != "" {
		return u.SpaceGUID
	}
	return u.OrgGUID
}

func (u *UsersInput) UniqueUsers() []string {
	return uniqueSlice(u.Users)
}

func (u *UsersInput) UniqueSPNUsers() []string {
	return uniqueSlice(u.SPNUsers)
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
		v := strings.Trim(strings.ToLower(value), " ")
		if _, ok := unique[v]; !ok {
			unique[v] = v
		}
	}
	for key := range unique {
		output = append(output, key)
	}
	return output
}
