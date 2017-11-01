package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gohttp "net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
)

type req struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type resp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestResponsesAndAdapters(t *testing.T) {
	cases := []struct {
		name     string
		verb     string
		path     string
		adapters []http.Adapter
		req      req
		want     response
		wantCode int
		wantErr  bool
	}{
		{
			name: "post handler",
			verb: "POST",
			path: "/svc/post_handler",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusOK,
				Data:   &resp{ID: 1, Name: "John Doe"},
				Errors: nil,
			},
			wantCode: gohttp.StatusOK,
			wantErr:  false,
		},
		{
			name: "post handler with apt",
			verb: "POST",
			path: "/svc/post_handler_apt",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusOK,
				Data:   &resp{ID: 1, Name: "msg1 msg2 John Doe"},
				Errors: nil,
			},
			adapters: []http.Adapter{
				func(h http.HandlerFunc) http.HandlerFunc {
					return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
						c = context.WithValue(c, "apt2", "msg2 ")
						h(c, w, r)
					}
				},
				func(h http.HandlerFunc) http.HandlerFunc {
					return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
						c = context.WithValue(c, "apt1", "msg1 ")
						h(c, w, r)
					}
				},
			},
			wantCode: gohttp.StatusOK,
			wantErr:  false,
		},
		{
			name: "endpoint handler",
			verb: "POST",
			path: "/svc/post_ep",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusOK,
				Data:   &resp{ID: 1, Name: "John Doe"},
				Errors: nil,
			},
			wantCode: gohttp.StatusOK,
			wantErr:  false,
		},
		{
			name: "endpoint handler with apt",
			verb: "POST",
			path: "/svc/post_ep_apt",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusOK,
				Data:   &resp{ID: 1, Name: "msg1 msg2 John Doe"},
				Errors: nil,
			},
			adapters: []http.Adapter{
				func(h http.HandlerFunc) http.HandlerFunc {
					return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
						c = context.WithValue(c, "apt2", "msg2 ")
						h(c, w, r)
					}
				},
				func(h http.HandlerFunc) http.HandlerFunc {
					return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
						c = context.WithValue(c, "apt1", "msg1 ")
						h(c, w, r)
					}
				},
			},
			wantCode: gohttp.StatusOK,
			wantErr:  false,
		},
		{
			name: "endpoint handler go error",
			verb: "POST",
			path: "/svc/post_ep_gerr",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusInternalServerError,
				Errors: []string{"endpoint error"},
			},
			wantCode: gohttp.StatusInternalServerError,
			wantErr:  false,
		},
		{
			name: "endpoint handler http error",
			verb: "POST",
			path: "/svc/post_ep_herr",
			req:  req{ID: 1, Name: "John Doe"},
			want: response{
				Code:   gohttp.StatusBadRequest,
				Errors: []string{"endpoint error"},
			},
			wantCode: gohttp.StatusBadRequest,
			wantErr:  false,
		},

		// TODO - Test adapters and errors
	}

	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			opts := []http.ServerOption{}
			if c.adapters != nil {
				opts = append(opts, http.WithAdapters(c.adapters...))
			}
			s := http.NewServer(opts...)

			svc := newHSvc()

			err := s.RegisterServices(svc)
			if c.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			go func(i int) {
				s.Run(8000 + i)
			}(i)
			defer s.Stop()

			body, _ := json.Marshal(c.req)
			req, _ := gohttp.NewRequest(c.verb, fmt.Sprintf("http://localhost:%d%s", 8000+i, c.path), bytes.NewReader(body))
			rsp, _ := gohttp.DefaultClient.Do(req)
			jresp := response{}
			json.NewDecoder(rsp.Body).Decode(&jresp)

			assert.Equal(t, c.want, jresp)
			assert.Equal(t, c.wantCode, rsp.StatusCode)
		})
	}
}

type response struct {
	Code   int      `json:"code"`
	Data   *resp    `json:"data,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func newHSvc() *hsvc {
	s := hsvc{}
	s.RegisterHandler("POST", "/post_handler", s.postHandler)
	s.RegisterHandler("POST", "/post_handler_apt", s.postHandlerApt)
	s.RegisterEndpoint("POST", "/post_ep", s.postEndpoint)
	s.RegisterEndpoint("POST", "/post_ep_apt", s.postEndpointApt)
	s.RegisterEndpoint("POST", "/post_ep_gerr", s.postEndpointGErr)
	s.RegisterEndpoint("POST", "/post_ep_herr", s.postEndpointHErr)
	return &s
}

type hsvc struct {
	http.BaseService
}

func (s *hsvc) Prefix() string { return "svc" }

func (s *hsvc) postHandler(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
	req := req{}
	json.NewDecoder(r.Body).Decode(&req)
	respond.WithJSON(w, r, http.NewResponse(resp{ID: req.ID, Name: req.Name}, gohttp.StatusOK))
}

func (s *hsvc) postHandlerApt(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
	req := req{}
	json.NewDecoder(r.Body).Decode(&req)
	msg1 := c.Value("apt1").(string)
	msg2 := c.Value("apt2").(string)
	respond.WithJSON(w, r, http.NewResponse(resp{ID: req.ID, Name: msg1 + msg2 + req.Name}, gohttp.StatusOK))
}

func (s *hsvc) postEndpoint(c context.Context, w gohttp.ResponseWriter, rq *req) (*http.Response, error) {
	return http.NewResponse(&resp{ID: rq.ID, Name: rq.Name}, gohttp.StatusOK), nil
}

func (s *hsvc) postEndpointApt(c context.Context, w gohttp.ResponseWriter, rq *req) (*http.Response, error) {
	msg1 := c.Value("apt1").(string)
	msg2 := c.Value("apt2").(string)
	return http.NewResponse(&resp{ID: rq.ID, Name: msg1 + msg2 + rq.Name}, gohttp.StatusOK), nil
}

func (s *hsvc) postEndpointGErr(c context.Context, w gohttp.ResponseWriter, rq *req) (*http.Response, error) {
	return nil, fmt.Errorf("endpoint error")
}

func (s *hsvc) postEndpointHErr(c context.Context, w gohttp.ResponseWriter, rq *req) (*http.Response, error) {
	return nil, http.WrapError(fmt.Errorf("endpoint error"), gohttp.StatusBadRequest)
}
