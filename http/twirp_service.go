package http

import (
	"context"
	"net/http"
	"strings"
)

// TwirpService is designed to be embedded into kit Service implementations
// in order to enable easier twirp Handler/hooks registration
type TwirpService struct {
	prefix       string
	twirpHandler http.Handler
}

// TwirpInit sets up kit endpoints and prefix for twirp handler
func (ts *TwirpService) TwirpInit(twirpPrefix string, twirpServer http.Handler) {
	ts.prefix = strings.Trim(twirpPrefix, "/")
	ts.twirpHandler = twirpServer
}

// Prefix returns service routing prefix
func (ts *TwirpService) Prefix() string { return ts.prefix }

// Endpoints returns all registered endpoints
func (ts *TwirpService) Endpoints() Endpoints {
	return Endpoints{
		"/{rest:.*}": &Endpoint{
			Methods: []string{"POST"},
			Handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				ctx = context.WithValue(ctx, ContextKey("HTTP-Authorization"), r.Header.Get("Authorization"))
				// TODO - Set other headers?
				ts.twirpHandler.ServeHTTP(w, r.WithContext(ctx))
			},
		},
	}
}
