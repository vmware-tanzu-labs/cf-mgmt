package http

import "net/http"

func IsResponseRedirect(statusCode int) bool {
	return IsStatusIn(statusCode,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther)
}

func IsStatusSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 399
}

func IsStatusIn(statusCode int, statuses ...int) bool {
	for _, s := range statuses {
		if statusCode == s {
			return true
		}
	}
	return false
}
