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
type Endpoints map[string]Endpoint
