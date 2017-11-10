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
	LDAPUsers          []string `long:"ldap-user" description:"Ldap User to add, specify multiple times"`
	LDAPUsersToRemove  []string `long:"ldap-user-to-remove" description:"Ldap User to remove, specify multiple times"`
	Users              []string `long:"user" description:"User to add, specify multiple times"`
	UsersToRemove      []string `long:"user-to-remove" description:"User to remove, specify multiple times"`
	SamlUsers          []string `long:"saml-user" description:"SAML user to add, specify multiple times"`
	SamlUsersToRemove  []string `long:"saml-user-to-remove" description:"SAML user to remove, specify multiple times"`
	LDAPGroups         []string `long:"ldap-group" description:"Group to add, specify multiple times"`
	LDAPGroupsToRemove []string `long:"ldap-group-to-remove" description:"Group to remove, specify multiple times"`
}

func updateUsersBasedOnRole(userMgmt *config.UserMgmt, currentLDAPGroups []string, userRole *UserRole, errorString *string) {
	userMgmt.LDAPGroups = removeFromSlice(addToSlice(currentLDAPGroups, userRole.LDAPGroups, errorString), userRole.LDAPGroupsToRemove)
	userMgmt.Users = removeFromSlice(addToSlice(userMgmt.Users, userRole.Users, errorString), userRole.UsersToRemove)
	userMgmt.SamlUsers = removeFromSlice(addToSlice(userMgmt.SamlUsers, userRole.SamlUsers, errorString), userRole.SamlUsersToRemove)
	userMgmt.LDAPUsers = removeFromSlice(addToSlice(userMgmt.LDAPUsers, userRole.LDAPUsers, errorString), userRole.LDAPUsersToRemove)
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

func addToSlice(theSlice, sliceToAdd []string, errorString *string) []string {
	checkForDuplicates(sliceToAdd, errorString)
	sliceToReturn := theSlice
	valuesThatExist := sliceToMap(theSlice)
	for _, val := range sliceToAdd {
		if _, ok := valuesThatExist[val]; !ok && val != "" {
			sliceToReturn = append(sliceToReturn, val)
		} else {
			*errorString += fmt.Sprintf("\n--value [%s] already exists in %v", val, theSlice)
		}
	}
	return sliceToReturn
}

func checkForDuplicates(slice []string, errorString *string) {
	sliceMap := make(map[string]string)
	for _, val := range slice {
		if _, ok := sliceMap[val]; ok && val != "" {
			*errorString += fmt.Sprintf("\n--value [%s] cannot be specified more than once %v", val, slice)
		} else {
			sliceMap[val] = val
		}
	}
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
