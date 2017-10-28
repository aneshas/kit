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
		adapter func(http.HandlerFunc) http.HandlerFunc
		aprime  func(http.HandlerFunc) http.HandlerFunc
		h       http.HandlerFunc
		want    string
	}{
		{
			name: "test post svc",
			verb: "POST",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "svc response")
			},
			want: "svc response",
		},
		{
			name: "test get svc",
			verb: "GET",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			},
			want: "get svc response",
		},
		{
			name: "test put svc with mw",
			verb: "PUT",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			},
			adapter: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					fmt.Fprintf(w, "hello from adapter ")
					h(c, w, r)
				}
			},
			want: "hello from adapter get svc response",
		},
		{
			name: "test put svc with stacked mw",
			verb: "PUT",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			},
			adapter: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					fmt.Fprintf(w, "hello from adapter ")
					h(c, w, r)
				}
			},
			aprime: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					fmt.Fprintf(w, "hello from adapter prime ")
					h(c, w, r)
				}
			},
			want: "hello from adapter prime hello from adapter get svc response",
		},
		{
			name: "test put svc with mw next",
			verb: "PUT",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				fmt.Fprintf(w, "get svc response")
			},
			adapter: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					h(c, w, r)
				}
			},
			want: "get svc response",
		},
		{
			name: "test mw context pass",
			verb: "PUT",
			path: "/svc",
			h: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
				v := c.Value("foo").(string)
				fmt.Fprintf(w, v)
			},
			adapter: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					c = context.WithValue(c, "foo", "foo msg")
					h(c, w, r)
				}
			},
			want: "foo msg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := svc{}
			s.RegisterHandler(c.verb, c.path, c.h)
			if c.adapter != nil {
				apts := []http.Adapter{}
				apts = append(apts, c.adapter)
				if c.aprime != nil {
					apts = append(apts, c.aprime)
				}
				s.Adapt(apts...)
			}
			endpoints := s.Endpoints()
			for path, ep := range endpoints {
				assert.Equal(t, c.path, path)
				assert.Equal(t, c.verb, ep.Methods[0])
				w := httptest.NewRecorder()
				ep.Handler(context.Background(), w, &gohttp.Request{})
				assert.Equal(t, c.want, string(w.Body.Bytes()))
			}

		})
	}
}

func TestRegisterEndpoint_WithEndpointsAndMW(t *testing.T) {
	cases := []struct {
		name     string
		verb     string
		path     string
		req      req
		endpoint interface{}
		adapter  func(http.HandlerFunc) http.HandlerFunc
		want     string
		wantErr  bool
		wantCode int
	}{
		{
			name: "test post svc",
			verb: "POST",
			path: "/svc",
			req:  req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				return http.NewResponse(resp{ID: req.ID, Name: req.Name}, gohttp.StatusOK), nil
			},
			want:     `{"code":200,"data":{"id":1,"name":"John"}}`,
			wantCode: gohttp.StatusOK,
		},
		{
			name: "test mw pass",
			verb: "POST",
			path: "/svc",
			req:  req{ID: 1, Name: "John"},
			adapter: func(h http.HandlerFunc) http.HandlerFunc {
				return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
					c = context.WithValue(c, "foo", "FooName")
					h(c, w, r)
				}
			},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				v := c.Value("foo").(string)
				return http.NewResponse(resp{ID: req.ID, Name: v}, gohttp.StatusOK), nil
			},
			want:     `{"code":200,"data":{"id":1,"name":"FooName"}}`,
			wantCode: gohttp.StatusOK,
		},
		{
			name:    "test invalid ctx",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c *req, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test invalid rw",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w *req, r *gohttp.Request, req *req) (*http.Response, error) {
				return nil, nil
			},
		},
		{
			name:    "test invalid r",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *req, req *req) (*http.Response, error) {
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
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) *http.Response {
				return nil
			},
		},
		{
			name:    "test invalid resp",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (error, *http.Response) {
				return nil, nil
			},
		},
		{
			name:    "test invalid err",
			verb:    "POST",
			path:    "/svc",
			wantErr: true,
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, *req) {
				return nil, nil
			},
		},
		{
			name: "test wrapped error",
			verb: "POST",
			path: "/svc",
			req:  req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				return nil, errors.Wrap(fmt.Errorf("endpoint error"), gohttp.StatusBadRequest)
			},
			want:     `{"code":400,"errors":["endpoint error"]}`,
			wantCode: gohttp.StatusBadRequest,
		},
		{
			name: "test go error",
			verb: "POST",
			path: "/svc",
			req:  req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				return nil, fmt.Errorf("endpoint error")
			},
			want:     `{"code":500,"errors":["endpoint error"]}`,
			wantCode: gohttp.StatusInternalServerError,
		},
		{
			name: "test nil response",
			verb: "POST",
			path: "/svc",
			req:  req{ID: 1, Name: "John"},
			endpoint: func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request, req *req) (*http.Response, error) {
				return nil, nil
			},
			want:     `{"code":200}`,
			wantCode: gohttp.StatusOK,
		},
		// TODO - Test Validation (Add Validate to *req) and err paths
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
				s.Adapt(c.adapter)
			}

			endpoints := s.Endpoints()

			for path, ep := range endpoints {
				assert.Equal(t, c.path, path)
				assert.Equal(t, c.verb, ep.Methods[0])

				w := httptest.NewRecorder()
				body, _ := json.Marshal(c.req)
				req, _ := gohttp.NewRequest(c.verb, "/svc", bytes.NewReader(body))
				ep.Handler(context.Background(), w, req)
				assert.Equal(t, c.want, string(w.Body.Bytes()))
			}
		})
	}
}
