package path

import (
	"fmt"
	"net/url"
	"strings"
)

func Format(urlFormat string, params ...any) string {
	// url encode any querystring params
	p := make([]any, len(params))
	for i, u := range params {
		switch v := u.(type) {
		case url.Values:
			p[i] = v.Encode()
		default:
			p[i] = u
		}
	}

	s := fmt.Sprintf(urlFormat, p...)
	return strings.TrimSuffix(s, "?")
}

func Join(elem ...string) string {
	var sb strings.Builder
	for _, e := range elem {
		if e == "" {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString("/")
			e = strings.TrimPrefix(e, "/")
		}
		e = strings.TrimSuffix(e, "/")
		sb.WriteString(e)
	}
	return sb.String()
}
