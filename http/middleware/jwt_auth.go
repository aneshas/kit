package middleware

import (
	"context"
	"net/http"
	"strings"
)

// WithTokenAuth represents http authentication middleware
// It uses auth func to check for valid session
func WithTokenAuth(af func(context.Context, string) error) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get("Authorization")
			if ah == "" {
				// respond.With(w, r, http.StatusBadRequest, fmt.Errorf("no authorization header found"))
				return
			}

			s := strings.Split(ah, " ")
			if len(s) < 2 || s[1] == "" {
				// respond.With(w, r, http.StatusBadRequest, fmt.Errorf("no bearer token found"))
				return
			}

			if err := af(context.Background(), s[1]); err != nil {
				// respond.With(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
