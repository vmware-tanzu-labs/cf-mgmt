package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type tokenPayload struct {
	Expiration int64 `json:"exp"`
}

func AccessTokenExpiration(accessToken string) (time.Time, error) {
	tp := strings.Split(accessToken, ".")
	if len(tp) != 3 {
		return time.Time{}, errors.New("access token format is invalid")
	}

	b := make([]byte, base64.RawURLEncoding.DecodedLen(len(tp[1])))
	_, err := base64.RawURLEncoding.Decode(b, []byte(tp[1]))
	if err != nil {
		return time.Time{}, errors.New("access token base64 encoding is invalid")
	}

	var t tokenPayload
	err = json.Unmarshal(b, &t)
	if err != nil {
		return time.Time{}, fmt.Errorf("access token is invalid: %w", err)
	}

	return time.Unix(t.Expiration, 0), nil
}
