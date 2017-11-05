package configcommands

import (
	"fmt"
	"strconv"

	"github.com/pivotalservices/cf-mgmt/config"
)

//BaseConfigCommand - commmand that specifies config-dir
type BaseConfigCommand struct {
	ConfigDirectory string `long:"config-dir" env:"CONFIG_DIR" default:"config" description:"Name of the config directory"`
}

type UserRole struct {
	LDAPUsers          []string `long:"ldap-user" description:"Ldap User to add, specify muliple times"`
	LDAPUsersToRemove  []string `long:"ldap-user-to-remove" description:"Ldap User to remove, specify muliple times"`
	Users              []string `long:"user" description:"User to add, specify muliple times"`
	UsersToRemove      []string `long:"user-to-remove" description:"User to remove, specify muliple times"`
	SamlUsers          []string `long:"saml-user" description:"SAML user to add, specify muliple times"`
	SamlUsersToRemove  []string `long:"saml-user-to-remove" description:"SAML user to remove, specify muliple times"`
	LDAPGroups         []string `long:"ldap-group" description:"User to add, specify muliple times"`
	LDAPGroupsToRemove []string `long:"ldap-group-to-remove" description:"User to remove, specify muliple times"`
}

func updateUsersBasedOnRole(userMgmt *config.UserMgmt, currentLDAPGroups []string, userRole *UserRole) {
	userMgmt.LDAPGroups = removeFromSlice(append(currentLDAPGroups, userRole.LDAPGroups...), userRole.LDAPGroupsToRemove)
	userMgmt.Users = removeFromSlice(append(userMgmt.Users, userRole.Users...), userRole.UsersToRemove)
	userMgmt.SamlUsers = removeFromSlice(append(userMgmt.SamlUsers, userRole.SamlUsers...), userRole.SamlUsersToRemove)
	userMgmt.LDAPUsers = removeFromSlice(append(userMgmt.LDAPUsers, userRole.LDAPUsers...), userRole.LDAPUsersToRemove)
	userMgmt.LDAPGroup = ""
}

func convertToInt(parameterName string, currentValue *int, proposedValue string, errorString *string) {
	if proposedValue == "" {
		return
	}
	i, err := strconv.Atoi(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an integer instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = i
}

func convertToBool(parameterName string, currentValue *bool, proposedValue string, errorString *string) {
	if proposedValue == "" {
		return
	}
	b, err := strconv.ParseBool(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an boolean instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = b
}

func removeFromSlice(theSlice, sliceToRemove []string) []string {
	var sliceToReturn []string
	valuesToRemove := sliceToMap(sliceToRemove)
	for _, val := range theSlice {
		if _, ok := valuesToRemove[val]; !ok && val != "" {
			sliceToReturn = append(sliceToReturn, val)
		}
	}
	return sliceToReturn
}

func sliceToMap(theSlice []string) map[string]string {
	theMap := make(map[string]string)
	for _, val := range theSlice {
		theMap[val] = val
	}
	return theMap
}
