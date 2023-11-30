package client

//go:generate go run ../tools/gen_error.go

import (
	"fmt"
)

type CloudFoundryHTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e CloudFoundryHTTPError) Error() string {
	return fmt.Sprintf("cfclient: HTTP error (%d): %s", e.StatusCode, e.Status)
}
