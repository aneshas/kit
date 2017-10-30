package adapter

import (
	"context"
	"fmt"
	gohttp "net/http"
	"strings"

	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

// AuthCallbackFunc represents auth callback that needs to be
// passed in to JWTAuth middleware. This func is called after token
// has been succesfully verified by adapter so client can do additional
// business auth check based on the token itself and claims extracted from it
// Actuall AuthCallbackFunc implementors should return error
// upon failed auth check or nil on success
type AuthCallbackFunc func(context.Context, string, map[string]interface{}) error 

// WithJWTAuth represents http authentication middleware
// It uses auth func to check for valid session
//
// - Alg
// - key
// - Token context key
func WithJWTAuth(callback AuthCallbackFunc) http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(ctx context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			ah := r.Header.Get("Authorization")
			if ah == "" {
				respond.WithJSON(
					w, r,
					http.WrapError(fmt.Errorf("no authorization header found"), gohttp.StatusBadRequest),
				)
				return
			}

			s := strings.Split(ah, " ")
			if len(s) < 2 || s[1] == "" {
				respond.WithJSON(
					w, r,
					http.WrapError(fmt.Errorf("no bearer token found"), gohttp.StatusBadRequest),
				)
				return
			}

			token := s[1]

			claims, err := verifyJWTToken(token)
			if err != nil {
				respond.WithJSON(w, r, err)
				return
			}

			if err := callback(context.Background(), token, claims); err != nil {
				respond.WithJSON(
					w, r,
					http.WrapError(fmt.Errorf("unauthorized: %v", err), gohttp.StatusUnauthorized),
				)
				return
			}

			h(context.WithValue(ctx, "auth-token", token), w, r)
		}
	}
}

func verifyJWTToken(token string) (map[string]interface{}, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// TODO - provide keyfile via parameter to New
		return []byte("Claire2016"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse provided token")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt token could not be verified")
}
