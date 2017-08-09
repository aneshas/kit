// Package respond provides common http response functionality
package respond

import (
	"encoding/json"
	"net/http"
)

var marshalError = `{"errors":["request was successfull but we were unable to encode the response."]}`

type response struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

// With makes a new json response based on a given response
// interface and code.
// Function checks for type of provided response interface
// and responds correctly based on the type of the interface value
func With(w http.ResponseWriter, r *http.Request, code int, i interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	var resp response

	if err, ok := i.(error); ok {
		resp.Errors = append(resp.Errors, err.Error())
	} else {
		resp.Data = i
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		b = []byte(marshalError)
		return
	}

	w.Write(b)
}
