// Package respond provides common http response functionality
package respond

import (
	"encoding/json"
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

	w.WriteHeader(gohttp.StatusOK)
	w.Write(
		encodeResp(
			response{
				Code: gohttp.StatusOK,
				Data: resp,
			},
		),
	)
}

func writeResponse(w gohttp.ResponseWriter, resp httpResponse) {
	w.WriteHeader(resp.Code())
	w.Write(
		encodeResp(
			response{
				Code: resp.Code(),
				Data: resp.Body(),
			},
		),
	)
}

func writeError(w gohttp.ResponseWriter, e httpError) {
	w.WriteHeader(e.Code())
	w.Write(
		encodeResp(
			response{
				Code:   e.Code(),
				Errors: []string{e.Err().Error()},
			},
		),
	)
}

func encodeResp(resp interface{}) []byte {
	data, err := json.Marshal(resp)
	if err != nil {
		return []byte(marshalError)
	}
	return data
}
