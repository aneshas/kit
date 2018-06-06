// Package http provides HTTP server implementation with server lifecycle control, commonly used
// adapters, easy service registration and response/error handling, tls support...
package http

import (
	"context"
	"net/http"
)

// ContextKey is a context key type
type ContextKey string

// Service defines http service interface
type Service interface {
	Prefix() string
	Endpoints() Endpoints
}

// Endpoint represents http api endpoint interface
type Endpoint struct {
	Methods []string
	Handler HandlerFunc
}

// HandlerFunc represents kit http handler func
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP implements http.Handler for HandlerFunc so it can be chained
// with third party middleware
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hf(context.Background(), w, r)
}

// Endpoints represents a map of service endpoints
type Endpoints map[string]*Endpoint

// NewResponse wraps provided code and resp into Response
// so it can be used with respond
func NewResponse(resp interface{}, code int) *Response {
	return &Response{
		code: code,
		body: resp,
	}
}

// Response represents http response
type Response struct {
	code int
	body interface{}
}

// Code returns response http code
func (r *Response) Code() int { return r.code }

// Body returns associated response body
func (r *Response) Body() interface{} { return r.body }

const contextReqKey = "tonto_http_request_key"

// ReqFromCtx returns http request associated with context
// (used with endpoints)
func ReqFromCtx(c context.Context) *http.Request {
	if r, ok := c.Value(contextReqKey).(*http.Request); ok {
		return r
	}
	return nil
}
