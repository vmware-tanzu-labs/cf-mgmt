package path

import (
	"errors"
	"net/url"
	"strconv"
)

type QuerystringReader struct {
	qs url.Values
}

func NewQuerystringReader(pageURL string) (*QuerystringReader, error) {
	if pageURL == "" {
		return nil, errors.New("cannot parse an empty pageURL")
	}
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	return &QuerystringReader{
		qs: u.Query(),
	}, nil
}

func (r QuerystringReader) String(key string) string {
	return r.qs.Get(key)
}

func (r QuerystringReader) Int(key string) int {
	i, _ := strconv.Atoi(r.qs.Get(key))
	return i
}
