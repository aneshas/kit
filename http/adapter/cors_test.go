package adapter_test

import (
	"context"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonto/kit/http/adapter"
	"github.com/tonto/kit/http/respond"
)

func TestWithCORS(t *testing.T) {
	cases := []struct {
		name        string
		origins     []string
		headers     []string
		methods     []string
		method      string
		host        string
		maxAge      int
		wantOrigins string
		wantMethods string
		wantHeaders string
		wantMaxAge  string
	}{
		{
			name:        "test asterisk origin",
			wantMaxAge:  "3600",
			origins:     []string{"*"},
			wantOrigins: "*",
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test origin",
			wantMaxAge:  "3600",
			origins:     []string{"foobar.com"},
			wantOrigins: "foobar.com",
			host:        "foobar.com",
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test origins",
			wantMaxAge:  "3600",
			origins:     []string{"foobar.com", "foobaz.org"},
			wantOrigins: "foobar.com",
			host:        "foobar.com",
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test asterisk origin override",
			wantMaxAge:  "3600",
			origins:     []string{"*", "foobaz.org"},
			wantOrigins: "*",
			host:        "foobar.com",
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test methods",
			wantMaxAge:  "3600",
			methods:     []string{"GET", "post", "PUT"},
			wantMethods: "GET, POST, PUT",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test headers",
			wantMaxAge:  "3600",
			headers:     []string{"Connection", "Content-Length"},
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Connection, Content-Length",
		},
		{
			name:        "test max age",
			maxAge:      86400,
			wantMaxAge:  "86400",
			wantMethods: "GET, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
			wantHeaders: "Accept, Accept-Language, Content-Language, Content-Type",
		},
		{
			name:        "test non options origin",
			method:      "GET",
			origins:     []string{"*"},
			wantOrigins: "*",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var opts []adapter.CORSOption

			opts = append(opts, adapter.WithCORSAllowOrigins(c.origins...))

			if c.methods != nil {
				opts = append(opts, adapter.WithCORSAllowMethods(c.methods...))
			}

			if c.headers != nil {
				opts = append(opts, adapter.WithCORSAllowHeaders(c.headers...))
			}

			if c.maxAge != 0 {
				opts = append(opts, adapter.WithCORSMaxAge(c.maxAge))
			}

			apt := adapter.WithCORS(opts...)

			hdlr := apt(func(ctx context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				respond.WithJSON(w, r, "response")
			})

			mtd := "OPTIONS"
			if c.method != "" {
				mtd = c.method
			}
			req, _ := gohttp.NewRequest(mtd, "/", nil)
			req.Host = c.host

			w := httptest.NewRecorder()
			hdlr(context.Background(), w, req)

			// resp := response{}
			// json.NewDecoder(w.Body).Decode(&resp)

			// assert.Equal(t, "response", resp.Data)
			assert.Equal(t, c.wantOrigins, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, c.wantMethods, w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, c.wantHeaders, w.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, c.wantMaxAge, w.Header().Get("Access-Control-Max-Age"))
		})
	}
}
