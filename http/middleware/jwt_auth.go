package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	respond "gopkg.in/matryer/respond.v1"
)

// Authorizer represents authorization domain service interface
type Authorizer interface {
	Authorize(context.Context, string) error
}

// WithJWTAuthMethod represents http authentication middleware
// It uses authorizer domain service to check for valid session
func WithJWTAuthMethod(j Authorizer) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get("Authorization")
			if ah == "" {
				respond.With(w, r, http.StatusBadRequest, fmt.Errorf("no authorization header found"))
				return
			}

			s := strings.Split(ah, " ")
			if len(s) < 2 || s[1] == "" {
				respond.With(w, r, http.StatusBadRequest, fmt.Errorf("no bearer token found"))
				return
			}

			if err := j.Authorize(context.Background(), s[1]); err != nil {
				respond.With(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
