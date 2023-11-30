package util

import "regexp"

// Matches - returns true if value matches any string in the list including regexes
func Matches(input string, list []string) bool {
	for _, name := range list {
		match, _ := regexp.MatchString(name, input)
		if match {
			return true
		}
	}
	return false
}
