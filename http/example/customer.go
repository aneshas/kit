package main

import (
	"context"
	"log"
	ghttp "net/http"

	"github.com/gorilla/mux"
	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
)

// NewCustomerService creates new customer service
func NewCustomerService(logger *log.Logger) *Customer {
	svc := Customer{
		logger: logger,
	}

	svc.Adapt(
		WithRequestLogger(logger),
	)

	svc.RegisterHandler("GET", "/details/{id}", svc.details)

	// If having problems accessing your endpoint you might
	// want to check err to see if you have endpoint validation issues
	err := svc.RegisterEndpoint("POST", "/create", svc.create)
	if err != nil {
		log.Fatal(err)
	}

	return &svc
}

// Customer represents customer http service
type Customer struct {
	http.BaseService

	logger *log.Logger
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
	return http.NewResponse(req, ghttp.StatusCreated), nil
}
