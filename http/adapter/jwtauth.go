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
// has been successfully verified by adapter so client can do additional
// business auth check based on the token itself and claims extracted from it
// Actuall AuthCallbackFunc implementors should return error
// upon failed auth check or nil on success
type AuthCallbackFunc func(context.Context, string, map[string]interface{}) error

// JWTAlg represents token signing alg type
type JWTAlg *jwt.SigningMethodHMAC

// JWTTokenKey is used to store token to context
const JWTTokenKey = "tonto_http_token_key"

var (
	// JWTAlgHS256 represents HMAC SHA256 token signing alg
	JWTAlgHS256 = jwt.SigningMethodHS256

	// JWTAlgHS384 represents HMAC SHA384 token signing alg
	JWTAlgHS384 = jwt.SigningMethodHS384

	// JWTAlgHS512 represents HMAC SHA512 token signing alg
	JWTAlgHS512 = jwt.SigningMethodHS512
)

// WithJWTAuth represents jwt authentication adapter
// It looks for bearer token in Authorization header, and if
// found tries to validate it against provided alg and key, if
// successful callback func is called to perform client side auth check.
func WithJWTAuth(alg JWTAlg, key []byte, callback AuthCallbackFunc) http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(ctx context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			ah := r.Header.Get("Authorization")
			if ah == "" {
				respond.WithJSON(
					w, r,
					http.NewError(gohttp.StatusBadRequest, fmt.Errorf("no authorization header found")),
				)
				return
			}

			s := strings.Split(ah, " ")
			if len(s) < 2 || s[1] == "" {
				respond.WithJSON(
					w, r,
					http.NewError(gohttp.StatusBadRequest, fmt.Errorf("no bearer token found")),
				)
				return
			}

			token := s[1]

			claims, err := verifyJWTToken(alg, token, key)
			if err != nil {
				respond.WithJSON(w, r, http.NewError(gohttp.StatusBadRequest, err))
				return
			}

			if err := callback(context.Background(), token, claims); err != nil {
				respond.WithJSON(
					w, r,
					http.NewError(gohttp.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err)),
				)
				return
			}

			h(context.WithValue(ctx, JWTTokenKey, token), w, r)
		}
	}
}

func verifyJWTToken(alg JWTAlg, token string, key []byte) (map[string]interface{}, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		if token.Method.Alg() != alg.Name {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse provided token")
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("jwt token could not be verified")
}
