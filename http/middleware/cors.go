package middleware

import "net/http"

// CORS creates a new CORS adapter
func CORS() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Add("Access-Control-Max-Age", "86400")
				w.WriteHeader(200)
				return
			}

			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		})
	}
}
