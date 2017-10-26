// Package http provides, commonly used http functionality such as:
// - Server lifecycle control (start, stop, status logging...)
// - Easy service registration and routing
// - Commonly used middleware implementations
package http

import (
	"net/http"
)

// Service defines http service interface
type Service interface {
	Prefix() string
	Endpoints() Endpoints
}

// Endpoint represents http api endpoint interface
type Endpoint struct {
	Methods []string
	Handler http.Handler
}

// Endpoints represents a map of service endpoints
type Endpoints map[string]*Endpoint

// NewResponse wraps provided code and resp into Response
// so it can be used with respond
func NewResponse(code int, resp interface{}) *Response {
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
