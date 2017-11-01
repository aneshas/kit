package main

import (
	"context"
	"log"
	ghttp "net/http"

	"github.com/tonto/kit/http"
)

// WithRequestLogger represents request logger adapter
func WithRequestLogger(logger *log.Logger) http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(c context.Context, w ghttp.ResponseWriter, r *ghttp.Request) {
			logger.Printf("IP: %s => %s %s", r.RemoteAddr, r.Method, r.URL.Path)
			h(c, w, r)
		}
	}
}
