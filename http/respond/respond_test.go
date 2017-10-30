package respond_test

import (
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
)

func TestWithJSON(t *testing.T) {
	cases := []struct {
		name       string
		resp       interface{}
		simpleResp string
		want       string
		wantCode   int
	}{
		{
			name: "test response success",
			resp: http.NewResponse(
				jresp{Foo: "val", Bar: 3},
				gohttp.StatusOK,
			),
			want:     `{"code":200,"data":{"foo":"val","bar":3}}`,
			wantCode: gohttp.StatusOK,
		},
		{
			name:       "test simple response success",
			simpleResp: "simple response",
			want:       `{"code":200,"data":"simple response"}`,
			wantCode:   gohttp.StatusOK,
		},
		{
			name: "test err success",
			resp: http.WrapError(
				fmt.Errorf("an error"),
				gohttp.StatusBadRequest,
			),
			want:     `{"code":400,"errors":["an error"]}`,
			wantCode: gohttp.StatusBadRequest,
		},
		{
			name: "test marshal err",
			resp: http.NewResponse(
				jresp{Foo: "error", Bar: 3},
				gohttp.StatusOK,
			),
			want:     `{"code":500,"errors":["request was successfull but we were unable to encode the response."]}`,
			wantCode: gohttp.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if c.resp != nil {
				respond.WithJSON(w, &gohttp.Request{}, c.resp)
			} else {
				respond.WithJSON(w, &gohttp.Request{}, c.simpleResp)
			}
			assert.Equal(t, c.want+"\n", string(w.Body.Bytes()))
			assert.Equal(t, c.wantCode, w.Code)
		})
	}
}

type response struct {
	Code   int      `json:"code"`
	Data   jresp    `json:"data,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

type jresp struct {
	Foo str `json:"foo"`
	Bar int `json:"bar"`
}

type str string

func (s str) MarshalJSON() ([]byte, error) {
	if s == "error" {
		return nil, fmt.Errorf("marshal error")
	}
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}
