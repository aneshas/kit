package adapter

import (
	"context"
	gohttp "net/http"

	"github.com/tonto/kit/http"
)

// WithCORS creates a new CORS adapter
func WithCORS() http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			if r.Method == "OPTIONS" {
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Add("Access-Control-Max-Age", "86400")
				w.WriteHeader(200)
				return
			}

			w.Header().Add("Access-Control-Allow-Origin", "*")
			h(c, w, r)
		}
	}
}
