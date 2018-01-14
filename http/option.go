package http

import (
	"log"
	"net/http"
)

// ServerOption is used for setting up
// server configuration
type ServerOption func(*Server)

// WithAdapters represents server option for setting up
// server-wide request adapters
func WithAdapters(a ...Adapter) ServerOption {
	return func(s *Server) {
		s.adapters = a
	}
}

// WithLogger represents server option for setting up logger
func WithLogger(l *log.Logger) ServerOption {
	return func(s *Server) {
		s.logger = l
	}
}

// WithTLSConfig represents server option for setting tls cer and key files
func WithTLSConfig(cert, key string) ServerOption {
	return func(s *Server) {
		s.certFile = cert
		s.keyFile = key
	}
}

// WithMux represents server option for setting a custom router.
// gorilla/mux is used as default
func WithMux(h http.Handler) ServerOption {
	return func(s *Server) {
		s.httpServer.Handler = h
	}
}

// WithNotFoundHandler represents server option for setting
// default not found handler
func WithNotFoundHandler(h http.Handler) ServerOption {
	return func(s *Server) {
		s.notFoundHandler = h
	}
}

// TODO - Add WithTLSConfig - so you can override the whole config
