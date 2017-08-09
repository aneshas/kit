package http

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/tonto/kit/http/respond"
)

type validator interface {
	Validate() error
}

type BaseService struct {
	endpoints Endpoints
}

// Prefix returns service routing prefix
func (b *BaseService) Prefix() string {
	return ""
}

// Endpoints returns all registered endpoints
func (b *BaseService) Endpoints() Endpoints {
	return b.endpoints
}

// RegisterEndpoint registeres service endpoint
func (b *BaseService) RegisterEndpoint(path string, h http.Handler, methods ...string) {
	if b.endpoints == nil {
		b.endpoints = make(map[string]Endpoint)
	}

	b.endpoints[path] = Endpoint{
		Methods: methods,
		Handler: h,
	}
}

// HandlerFromMethod creates new handler from a given service method.
// Required request struct will be recognised and request body will be
// correctly unmarshaled to it.
func (b *BaseService) HandlerFromMethod(m interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := reflect.ValueOf(m)
		reqParamType := v.Type().In(2).Elem()
		req := reflect.New(reqParamType).Interface()

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		if validator, ok := interface{}(req).(validator); ok {
			err = validator.Validate()
			if err != nil {
				respond.With(w, r, http.StatusBadRequest, err)
				return
			}
		}

		v.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			reflect.ValueOf(req),
		})
	})
}
