package middleware

import "net/http"

// Adapter represents http.Handler adapter type
type Adapter func(http.Handler) http.Handler

// Adapt decorates given Handler with provided adapters
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// Adapt decorates given HandlerFunc with provided adapters
// and returns a new Handler
func AdaptFunc(h http.HandlerFunc, adapters ...Adapter) http.Handler {
	var handler http.Handler
	for _, adapter := range adapters {
		handler = adapter(http.HandlerFunc(h))
	}
	return handler
}
