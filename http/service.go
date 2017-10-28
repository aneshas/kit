package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/tonto/kit/http/respond"
)

type validator interface {
	Validate() error
}

// BaseService represents base http service
type BaseService struct {
	m         sync.Mutex
	endpoints Endpoints
	mw        []Adapter
}

// Prefix returns service routing prefix
func (b *BaseService) Prefix() string { return "/" }

// RegisterHandler is a helper method that registers service HandlerFunc
// Service HandlerFunc is an extension of http.HandlerFunc which only adds context.Context
// as first parameter, the rest stays the same
func (b *BaseService) RegisterHandler(verb string, path string, h HandlerFunc) {
	if b.endpoints == nil {
		b.endpoints = make(map[string]*Endpoint)
	}
	b.endpoints[path] = &Endpoint{
		Methods: []string{verb},
		Handler: h,
	}
}

// RegisterEndpoint is a helper method that registers service json endpoint
// JSON endpoint method should have the following signature:
// func(c context.Context, w http.ResponseWriter, r *http.Request, req *CustomeType) (*http.Response, error)
// where *CustomType is your custom request type to which r.Body will be json unmarshalled automatically
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

func (b *BaseService) handlerFromMethod(m interface{}) (HandlerFunc, error) {
	err := b.checkMtdSig(m)
	if err != nil {
		return nil, err
	}

	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		req, err := b.decodeReq(r, m)
		if err != nil {
			respond.WithJSON(
				w, r,
				WrapError(fmt.Errorf("internal error: could not decode request"), http.StatusBadRequest),
			)
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
			reflect.ValueOf(c),
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			reflect.ValueOf(req),
		})

		if !ret[1].IsNil() {
			b.writeError(w, r, ret[1].Interface())
			return
		}

		if ret[0].IsNil() {
			respond.WithJSON(w, r, NewResponse(nil, http.StatusOK))
			return
		}

		resp := ret[0].Interface().(*Response)
		respond.WithJSON(w, r, resp)
	}, nil
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
		return nil, fmt.Errorf("error decoding json")
	}

	return req, nil
}

func (b *BaseService) writeError(w http.ResponseWriter, r *http.Request, e interface{}) {
	if _, ok := e.(*Error); ok {
		respond.WithJSON(w, r, e)
		return
	}
	respond.WithJSON(w, r, WrapError(e.(error), http.StatusInternalServerError))
}

// Endpoints returns all registered endpoints
func (b *BaseService) Endpoints() Endpoints {
	for _, e := range b.endpoints {
		if b.mw != nil {
			e.Handler = AdaptHandlerFunc(e.Handler, b.mw...)
		}
	}
	return b.endpoints
}

// Adapt is used to adapt the service with provided adapters
func (b *BaseService) Adapt(mw ...Adapter) { b.mw = mw }
