package http

import (
	"context"
	"io"
)

const DefaultContentType = "application/json"

// Request is used to help build up an HTTP request
type Request struct {
	method          string
	pathAndQuery    string
	contentType     string
	contentLength   *int64
	followRedirects bool
	context         context.Context

	// can set one or the other but not both
	body   io.Reader
	object any

	// arbitrary headers
	headers map[string]string
}

// NewRequest creates a new minimally configured HTTP request instance
func NewRequest(ctx context.Context, method, pathAndQuery string) *Request {
	return &Request{
		context:         ctx,
		method:          method,
		pathAndQuery:    pathAndQuery,
		headers:         make(map[string]string),
		followRedirects: true,
	}
}

// WithObject adds an object to the request body to be JSON serialized
func (r *Request) WithObject(obj any) *Request {
	r.object = obj

	// default content type to json if provided an object
	if r.contentType == "" {
		r.contentType = DefaultContentType
	}
	return r
}

// WithBody adds the specified body as-is to the request and defaults the content type to JSON
func (r *Request) WithBody(body io.Reader) *Request {
	r.body = body

	// default content type to json if provided a body
	if r.contentType == "" {
		r.contentType = DefaultContentType
	}
	return r
}

// WithContentType sets the content type of the request body
func (r *Request) WithContentType(contentType string) *Request {
	r.contentType = contentType
	return r
}

// WithContentLength sets the content length, needed for file uploads
func (r *Request) WithContentLength(len int64) *Request {
	r.contentLength = &len
	return r
}

// WithHeader sets an arbitrary header on the request
func (r *Request) WithHeader(name, value string) *Request {
	r.headers[name] = value
	return r
}

// WithFollowRedirects sets the content type of the request body
func (r *Request) WithFollowRedirects(follow bool) *Request {
	r.followRedirects = follow
	return r
}
