package util

import (
	"net/http"
	"net/http/httputil"

	"github.com/xchapter7x/lo"
)

type LoggingTransport struct {
	base http.RoundTripper
}

func NewLoggingTransport(roundTripper http.RoundTripper) *LoggingTransport {
	return &LoggingTransport{base: roundTripper}
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logRequest(req)

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *LoggingTransport) logRequest(req *http.Request) {
	bytes, _ := httputil.DumpRequest(req, true)
	lo.G.Infof("Request: [%s]", string(bytes))
}

func (t *LoggingTransport) logResponse(resp *http.Response) {
	bytes, _ := httputil.DumpResponse(resp, true)
	lo.G.Infof("Response: [%s]", string(bytes))
}
