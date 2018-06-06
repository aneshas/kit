package main

import (
	"context"
	"fmt"
	ghttp "net/http"

	"github.com/gorilla/mux"
	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
)

// NewCustomerService creates new customer service
func NewCustomerService(apts ...http.Adapter) *Customer {
	svc := Customer{}

	svc.RegisterHandler("GET", "/details/{id}", svc.details)
	svc.RegisterEndpoint("PUT", "/update", svc.update)

	svc.MustRegisterEndpoint("POST", "/create", svc.create)

	svc.Adapt(apts...)

	return &svc
}

// Customer represents customer http service
type Customer struct {
	http.BaseService
}

// Prefix returns service routing prefix
func (o *Customer) Prefix() string { return "customer" }

func (o *Customer) details(c context.Context, w ghttp.ResponseWriter, r *ghttp.Request) {
	id := mux.Vars(r)["id"]
	respond.WithJSON(w, r, http.NewResponse(id, ghttp.StatusOK))
}

type customerCreateReq struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (o *Customer) create(c context.Context, w ghttp.ResponseWriter, req *customerCreateReq) (*http.Response, error) {
	return nil, http.NewError(
		ghttp.StatusInternalServerError,
		fmt.Errorf("some error"),
		fmt.Errorf("another error"),
	)
}

type customerUpdateReq struct{}

func (o *Customer) update(c context.Context, w ghttp.ResponseWriter, req *customerUpdateReq) (*http.Response, error) {
	// We can also return plain error in which case default status is 500
	return nil, fmt.Errorf("internal error")
}
