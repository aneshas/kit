// Package respond provides common http response functionality
package respond

import (
	"encoding/json"
	"io"
	gohttp "net/http"
)

var marshalError = `{"code":500,"errors":["request was successfull but we were unable to encode the response."]}`

type response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

type httpResponse interface {
	Code() int
	Body() interface{}
}

type httpError interface {
	Code() int
	Err() error
}

// WithJSON makes a new json response based on a given response interface
// If provided resp is of type errors.Error error response will be made,
// otherwise provider resp will be json encoded and written to w
func WithJSON(w gohttp.ResponseWriter, r *gohttp.Request, resp interface{}) {
	w.Header().Add("Content-Type", "application/json")

	hresp, ok := resp.(httpResponse)
	if ok {
		writeResponse(w, hresp)
		return
	}

	herr, ok := resp.(httpError)
	if ok {
		writeError(w, herr)
		return
	}

	err, ok := resp.(error)
	if ok {
		writeSimpleError(w, err)
		return
	}

	w.WriteHeader(gohttp.StatusOK)
	writeJSON(
		w,
		response{
			Code: gohttp.StatusOK,
			Data: resp,
		},
	)
}

func writeResponse(w gohttp.ResponseWriter, resp httpResponse) {
	w.WriteHeader(resp.Code())
	writeJSON(
		w,
		response{
			Code: resp.Code(),
			Data: resp.Body(),
		},
	)
}

func writeError(w gohttp.ResponseWriter, e httpError) {
	w.WriteHeader(e.Code())
	writeJSON(
		w,
		response{
			Code:   e.Code(),
			Errors: []string{e.Err().Error()},
		},
	)
}

func writeSimpleError(w gohttp.ResponseWriter, err error) {
	w.WriteHeader(gohttp.StatusInternalServerError)
	writeJSON(
		w,
		response{
			Code:   gohttp.StatusInternalServerError,
			Errors: []string{err.Error()},
		},
	)
}

func writeJSON(w io.Writer, resp interface{}) {
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.Write([]byte(marshalError + "\n"))
	}
}
