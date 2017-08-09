package http

import (
	"log"

	"github.com/tonto/kit/http/middleware"
)

// ServerOption is used for setting up
// server configuration
type ServerOption func(*HTTPServer)

// WithMiddleware represents server option for setting up
// pre request middlewares
func WithMiddleware(m ...middleware.Adapter) ServerOption {
	return func(s *HTTPServer) {
		s.httpSrvr.Handler = middleware.Adapt(s.router, m...)
	}
}

// WithLogger is used for setting up server logger
func WithLogger(l *log.Logger) ServerOption {
	return func(s *HTTPServer) {
		s.logger = l
	}
}
