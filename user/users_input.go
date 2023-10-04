package user

// UsersInput
type UsersInput struct {
	SpaceGUID                                         string
	OrgGUID                                           string
	LdapUsers, Users, SPNUsers, GroupNames, SamlUsers []string
	SpaceName                                         string
	OrgName                                           string
	RemoveUsers                                       bool
	RoleUsers                                         *RoleUsers
	AddUser                                           func(updateUserInput UsersInput, userName, userGUID string) error
	RemoveUser                                        func(updateUserInput UsersInput, userName, userGUID string) error
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

func (u *UsersInput) UniqueGroupNames() []string {
	return uniqueSlice(u.GroupNames)
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
