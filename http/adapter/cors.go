package adapter

import (
	"context"
	"fmt"
	gohttp "net/http"
	"strings"

	"github.com/tonto/kit/http"
)

// CORSOption represents cors option
type CORSOption func(*corsCfg)

// WithCORS creates a new CORS adapter
func WithCORS(opts ...CORSOption) http.Adapter {
	cfg := corsCfg{}
	for _, o := range opts {
		o(&cfg)
	}
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			if cfg.origins != nil {
				var origins string
				for _, o := range cfg.origins {
					if o == "*" {
						origins = "*"
						break
					}
					if r.Host == o {
						origins = o
						break
					}
				}
				if origins != "" {
					w.Header().Add("Access-Control-Allow-Origin", origins)
				}
			}

			if r.Method == "OPTIONS" {
				methods := cfg.methods
				if methods == "" {
					methods = "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH"
				}
				w.Header().Add("Access-Control-Allow-Methods", methods)

				headers := cfg.headers
				if headers == "" {
					headers = "Authorization, Accept, Accept-Language, Content-Language, Content-Type"
				}
				w.Header().Add("Access-Control-Allow-Headers", headers)

				age := cfg.maxAge
				if age == "" {
					age = "3600"
				}
				w.Header().Add("Access-Control-Max-Age", age)

				w.Header().Add("Access-Control-Allow-Credentials", "true")

				w.WriteHeader(200)
				return
			}

			h(c, w, r)
		}
	}
}

type corsCfg struct {
	origins []string
	methods string
	headers string
	maxAge  string
}

// WithCORSAllowOrigins sets allowed origins
func WithCORSAllowOrigins(origins ...string) CORSOption {
	return func(cfg *corsCfg) {
		cfg.origins = origins
	}
}

// WithCORSAllowMethods sets allowed methods
func WithCORSAllowMethods(methods ...string) CORSOption {
	return func(cfg *corsCfg) {
		for _, mtd := range methods {
			cfg.methods += strings.ToUpper(mtd) + ", "
		}
		cfg.methods = strings.TrimRight(cfg.methods, ", ")
	}
}

// WithCORSAllowHeaders sets allowed headers
func WithCORSAllowHeaders(headers ...string) CORSOption {
	return func(cfg *corsCfg) {
		for _, h := range headers {
			cfg.headers += h + ", "
		}
		cfg.headers = strings.TrimRight(cfg.headers, ", ")
	}
}

// WithCORSMaxAge sets cors max age
func WithCORSMaxAge(age int) CORSOption {
	return func(cfg *corsCfg) {
		cfg.maxAge = fmt.Sprintf("%d", age)
	}
}
