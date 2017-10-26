package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/errors"
)

type svc struct {
	http.BaseService
}

func TestDefaultPrefix(t *testing.T) {
	s := svc{}
	assert.Equal(t, "/", s.Prefix())
}

func TestRegisterHandler_WithEndpointsAndMW(t *testing.T) {
	cases := []struct {
		name    string
		verb    string
		path    string
		adapter func(gohttp.Handler) gohttp.Handler
		h       gohttp.Handler
		want    string
	}{
		{
			name: "test post svc",
			verb: "POST",
			path: "/svc",
			h: gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "svc response")
			}),
			want: "svc response",
		},
		{
			name: "test get svc",
			verb: "GET",
			path: "/svc",
			h: gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			}),
			want: "get svc response",
		},
		{
			name: "test put svc with mw",
			verb: "PUT",
			path: "/svc",
			h: gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			}),
			adapter: func(h gohttp.Handler) gohttp.Handler {
				return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
					fmt.Fprintf(w, "hello from adapter ")
					h.ServeHTTP(w, r)
				})
			},
			want: "hello from adapter get svc response",
		},
		{
			name: "test put svc with mw next",
			verb: "PUT",
			path: "/svc",
			h: gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			}),
			adapter: func(h gohttp.Handler) gohttp.Handler {
				return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
					h.ServeHTTP(w, r)
				})
			},
			want: "get svc response",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := svc{}
			s.RegisterHandler(c.verb, c.path, c.h)
			if c.adapter != nil {
				s.RegisterMiddleware(c.adapter)
			}
			endpoints := s.Endpoints()
			for path, ep := range endpoints {
				assert.Equal(t, c.path, path)
				assert.Equal(t, c.verb, ep.Methods[0])
				w := httptest.NewRecorder()
				ep.Handler.ServeHTTP(w, &gohttp.Request{})
				assert.Equal(t, c.want, string(w.Body.Bytes()))
			}

		})
	}
}

func TestRegisterEndpoint_WithEndpointsAndMW(t *testing.T) {
	type Req struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type Resp struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	cases := []struct {
		name     string
		verb     string
		path     string
		req      Req
		endpoint interface{}
		adapter  func(gohttp.Handler) gohttp.Handler
		want     string
		wantErr  bool
		wantCode int
	}{
		{
			name: "test post svc",
			verb: "POST",
			path: "/svc",
			req:  Req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, error) {
				return http.NewResponse(gohttp.StatusOK, Resp{ID: req.ID, Name: req.Name}), nil
			},
			want:     `{"code":200,"data":{"id":1,"name":"John"}}`,
			wantCode: gohttp.StatusOK,
		},
		{
			name:    "test invalid ctx",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c *Req, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test invalid rw",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w *Req, r *gohttp.Request, req *Req) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test invalid r",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *Req, req *Req) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test num params",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test num ret",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) *http.Response {
				return nil
			},
		},
		{
			name:    "test invalid resp",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (error, *http.Response) {
				return nil, nil
			},
		},
		{
			name:    "test invalid err",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, *Req) {
				return nil, nil
			},
		},
		{
			name: "test wrapped error",
			verb: "POST",
			path: "/svc",
			req:  Req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, error) {
				return nil, errors.Wrap(fmt.Errorf("endpoint error"), gohttp.StatusBadRequest)
			},
			want:     `{"code":400,"errors":["endpoint error"]}`,
			wantCode: gohttp.StatusBadRequest,
		},
		{
			name: "test go error",
			verb: "POST",
			path: "/svc",
			req:  Req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, error) {
				return nil, fmt.Errorf("endpoint error")
			},
			want:     `{"code":500,"errors":["endpoint error"]}`,
			wantCode: gohttp.StatusInternalServerError,
		},
		{
			name: "test nil response",
			verb: "POST",
			path: "/svc",
			req:  Req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *Req) (*http.Response, error) {
				return nil, nil
			},
			want:     `{"code":200}`,
			wantCode: gohttp.StatusOK,
		},
		// TODO - Test Validation (Add Validate to *Req)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := svc{}

			err := s.RegisterEndpoint(c.verb, c.path, c.endpoint)
			if c.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			if c.adapter != nil {
				s.RegisterMiddleware(c.adapter)
			}

			endpoints := s.Endpoints()

			for path, ep := range endpoints {
				assert.Equal(t, c.path, path)
				assert.Equal(t, c.verb, ep.Methods[0])

				w := httptest.NewRecorder()
				body, _ := json.Marshal(c.req)
				req, _ := gohttp.NewRequest(c.verb, "/svc", bytes.NewReader(body))
				ep.Handler.ServeHTTP(w, req)
				assert.Equal(t, c.want, string(w.Body.Bytes()))
			}
		})
	}
}
