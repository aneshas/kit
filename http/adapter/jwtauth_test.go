package adapter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonto/kit/http/adapter"
	"github.com/tonto/kit/http/respond"
)

func TestWithJWTAuth(t *testing.T) {
	cases := []struct {
		name     string
		alg      adapter.JWTAlg
		token    string
		header   string
		authErr  error
		tokenKey string
		claims   map[string]string
		key      []byte
		want     response
	}{
		{
			name:     "test HS256",
			alg:      adapter.JWTAlgHS256,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.RQX-U1ElsPmUW__sZbJjhPOG6G8F0hYUnKNlE1bGR9k",
			want: response{
				Code: 200,
				Data: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.RQX-U1ElsPmUW__sZbJjhPOG6G8F0hYUnKNlE1bGR9k",
			},
		},
		{
			name:     "test HS384",
			alg:      adapter.JWTAlgHS384,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzM4NCJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.jCSHnlbSNT_JdiJX-Ue9TGCFuwBoru3yOAWDNk5wApdJQigZMst0xjCzc0QEBlsq",
			want: response{
				Code: 200,
				Data: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzM4NCJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.jCSHnlbSNT_JdiJX-Ue9TGCFuwBoru3yOAWDNk5wApdJQigZMst0xjCzc0QEBlsq",
			},
		},
		{
			name:     "test HS512",
			alg:      adapter.JWTAlgHS512,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			want: response{
				Code: 200,
				Data: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			},
		},
		{
			name:     "test HS512",
			alg:      adapter.JWTAlgHS512,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			claims: map[string]string{
				"Surname": "Rocket",
				"Email":   "jrocket@example.com",
			},
			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			want: response{
				Code: 200,
				Data: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			},
		},
		{
			name:     "test no token",
			alg:      adapter.JWTAlgHS512,
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "",
			want: response{
				Code:   400,
				Errors: []string{"no authorization header found"},
			},
		},
		{
			name:     "test no token",
			alg:      adapter.JWTAlgHS512,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "",
			want: response{
				Code:   400,
				Errors: []string{"no bearer token found"},
			},
		},
		{
			name:     "test unauthorized",
			alg:      adapter.JWTAlgHS512,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			authErr:  fmt.Errorf("auth error"),
			want: response{
				Code:   401,
				Errors: []string{"unauthorized: auth error"},
			},
		},
		{
			name:     "test token parse error",
			alg:      adapter.JWTAlgHS256,
			header:   "Authorization",
			tokenKey: "tkey",
			key:      []byte("123456"),
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1MDkzNzczMDEsImV4cCI6MTU0MDkxMzMwMSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.pQnK-DKhBGOMig8dDvQztdWkKl51mhvJeZujoHjAoXCYFPv6UJlw19RlCczoqmqqsK2fAjYnDgUiYGDSvhISmw",
			want: response{
				Code:   400,
				Errors: []string{"could not parse provided token"},
			},
		},

		// TODO - Test claims validation
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			apt := adapter.WithJWTAuth(
				c.alg,
				c.key,
				c.tokenKey,
				func(ctx context.Context, token string, claims map[string]interface{}) error {
					log.Println(claims)
					c := c
					if c.claims != nil {
						for claim, val := range c.claims {
							jval, ok := claims[claim]
							if !ok {
								t.Fail()
								return c.authErr
							}
							assert.Equal(t, val, jval.(string))
						}
					}
					return c.authErr
				},
			)

			hdlr := apt(func(ctx context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				respond.WithJSON(w, r, ctx.Value(c.tokenKey).(string))
			})

			req, _ := gohttp.NewRequest("GET", "/", nil)
			req.Header.Add(c.header, "Bearer "+c.token)

			w := httptest.NewRecorder()
			hdlr(context.Background(), w, req)

			resp := response{}
			json.NewDecoder(w.Body).Decode(&resp)

			assert.Equal(t, c.want, resp)
		})
	}
}

type response struct {
	Code   int      `json:"code"`
	Data   string   `json:"data"`
	Errors []string `json:"errors"`
}
