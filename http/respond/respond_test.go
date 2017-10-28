package respond_test

import (
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/errors"
	"github.com/tonto/kit/http/respond"
)

func TestWithJSON(t *testing.T) {
	cases := []struct {
		name     string
		resp     interface{}
		want     string
		wantCode int
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
			name: "test err success",
			resp: errors.Wrap(
				fmt.Errorf("an error"),
				gohttp.StatusBadRequest,
			),
			want:     `{"code":400,"errors":["an error"]}`,
			wantCode: gohttp.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respond.WithJSON(w, &gohttp.Request{}, c.resp)
			assert.Equal(t, c.want, string(w.Body.Bytes()))
			assert.Equal(t, c.wantCode, w.Code)
		})
	}
}

type jresp struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
