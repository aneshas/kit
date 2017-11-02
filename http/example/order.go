package main

import (
	"context"
	"fmt"
	"log"
	ghttp "net/http"

	"github.com/gorilla/mux"

	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/respond"
)

// NewOrderService creates new order service
func NewOrderService(logger *log.Logger) *Order {
	svc := Order{
		logger: logger,
	}

	svc.Adapt(
		WithRequestLogger(logger),
	)

	// Normal handler where you handle request decoding and validation
	svc.RegisterHandler("GET", "/details{id}", svc.details)

	// Register json endpoint which handles automatic validation, request
	// decoding and easy to use return based response
	svc.RegisterEndpoint("POST", "/create", svc.create)

	return &svc
}

// Order represents order http service
type Order struct {
	http.BaseService

	logger *log.Logger
}

// Prefix returns service routing prefix
func (o *Order) Prefix() string { return "order" }

func (o *Order) details(c context.Context, w ghttp.ResponseWriter, r *ghttp.Request) {
	id := mux.Vars(r)["id"]
	respond.WithJSON(w, r, http.NewResponse(id, ghttp.StatusOK))
}

type orderCreateReq struct {
	CustomerID int64 `json:"customer_id"`
}

// Implement http.Validator to make use of automatic request validation
func (r *orderCreateReq) Validate() error {
	if r.CustomerID == 0 {
		return fmt.Errorf("customer id must not be empty")
	}
	return nil
}

type orderCreateResp struct {
	CustomerID int64 `json:"customer_id"`
}

func (o *Order) create(c context.Context, w ghttp.ResponseWriter, req *orderCreateReq) (*http.Response, error) {
	return http.NewResponse(
		orderCreateResp{
			// Return customer id to demonstrate that
			// we got it in the request
			CustomerID: req.CustomerID,
		},
		ghttp.StatusCreated,
	), nil
}
