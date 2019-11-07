package organization

import (
	"regexp"
)

func Matches(orgName string, orgList []string) bool {
	for _, name := range orgList {
		match, _ := regexp.MatchString(name, orgName)
		if match {
			return true
		}
	}
	return false
}
