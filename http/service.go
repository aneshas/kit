package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/tonto/kit/http/errors"
	"github.com/tonto/kit/http/middleware"
	"github.com/tonto/kit/http/respond"
)

type validator interface {
	Validate() error
}

// BaseService represents base http service
type BaseService struct {
	endpoints Endpoints
	mw        []middleware.Adapter
}

// Prefix returns service routing prefix
func (b *BaseService) Prefix() string { return "/" }

// RegisterHandler is a helper method that registers service endpoint handler
func (b *BaseService) RegisterHandler(verb string, path string, h http.Handler) {
	if b.endpoints == nil {
		b.endpoints = make(map[string]*Endpoint)
	}
	b.endpoints[path] = &Endpoint{
		Methods: []string{verb},
		Handler: h,
	}
}

// RegisterEndpoint is a helper method that registers service endpoint
// method should have a following signature:
// func(c context.Context, w http.ResponseWriter, r *http.Request, req *CustomeType) (*http.Response, error)
func (b *BaseService) RegisterEndpoint(verb string, path string, method interface{}) error {
	h, err := b.handlerFromMethod(method)
	if err != nil {
		return err
	}

	if b.endpoints == nil {
		b.endpoints = make(map[string]*Endpoint)
	}

	b.endpoints[path] = &Endpoint{
		Methods: []string{verb},
		Handler: h,
	}

	return nil
}

func (b *BaseService) handlerFromMethod(m interface{}) (http.Handler, error) {
	err := b.checkMtdSig(m)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := b.decodeReq(r, m)
		if err != nil {
			fmt.Fprintf(w, "internal error - could not decode request")
			return
		}

		if validator, ok := interface{}(req).(validator); ok {
			err = validator.Validate()
			if err != nil {
				// respond.With(w, r, http.StatusBadRequest, err)
				return
			}
		}

		v := reflect.ValueOf(m)
		ret := v.Call([]reflect.Value{
			reflect.ValueOf(context.Background()),
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			reflect.ValueOf(req),
		})

		if !ret[1].IsNil() {
			b.writeError(w, r, ret[1].Interface())
			return
		}

		if ret[0].IsNil() {
			respond.WithJSON(w, r, NewResponse(http.StatusOK, nil))
			return
		}

		resp := ret[0].Interface().(*Response)
		respond.WithJSON(w, r, resp)
	}), nil
}

func (b *BaseService) checkMtdSig(m interface{}) error {
	t := reflect.ValueOf(m).Type()

	if t.NumIn() != 4 {
		return fmt.Errorf("incorrect endpoint signature (must have 4 params - refer to docs)")
	}

	if t.NumOut() != 2 {
		return fmt.Errorf("incorrect endpoint signature (must have 2 ret vals - refer to docs)")
	}

	if !t.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return fmt.Errorf("param one must implement context.Context")
	}

	if !t.In(1).Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
		return fmt.Errorf("param two must implement http.ResponseWriter")
	}

	if t.In(2) != reflect.TypeOf(&http.Request{}) {
		return fmt.Errorf("param three must be of type *http.Request")
	}

	if t.Out(0) != reflect.TypeOf(&Response{}) {
		return fmt.Errorf("first ret value must be of type *kit/http/Response")
	}

	if !t.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("second ret value must implement error interface")
	}

	return nil
}

func (b *BaseService) decodeReq(r *http.Request, m interface{}) (interface{}, error) {
	defer r.Body.Close()

	v := reflect.ValueOf(m)
	reqParamType := v.Type().In(3).Elem()
	req := reflect.New(reqParamType).Interface()

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		// respond.With(w, r, http.StatusBadRequest, err)
		return nil, fmt.Errorf("error decoding json")
	}

	return req, nil
}

func (b *BaseService) writeError(w http.ResponseWriter, r *http.Request, e interface{}) {
	if _, ok := e.(*errors.Error); ok {
		respond.WithJSON(w, r, e)
		return
	}
	respond.WithJSON(w, r, errors.Wrap(e.(error), http.StatusInternalServerError))
}

// Endpoints returns all registered endpoints
func (b *BaseService) Endpoints() Endpoints {
	for _, e := range b.endpoints {
		if b.mw != nil {
			e.Handler = middleware.Adapt(e.Handler, b.mw...)
		}
	}
	return b.endpoints
}

// RegisterMiddleware is a helper method that registers provided middlewares
// for service wide usage, ie. provided middlewares are applied to all endpoints
func (b *BaseService) RegisterMiddleware(mw ...middleware.Adapter) { b.mw = mw }
