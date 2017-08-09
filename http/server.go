// Package server provides common http server functionality
package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

var (
	ErrNoEndpointsDefined = errors.New("service has no endpoints defined")
	ErrNoServices         = errors.New("no services provided")
)

// NewServer creates new http server
// TODO - Add options eg. WithPort...
// WithLogger etc...
func NewServer(opts ...ServerOption) *HTTPServer {
	srvr := HTTPServer{
		httpSrvr: &http.Server{},
		router:   mux.NewRouter(),
	}

	for _, o := range opts {
		o(&srvr)
	}

	if srvr.logger == nil {
		srvr.logger = log.New(os.Stdout, "http ", log.Ldate|log.Ltime|log.Llongfile)
	}

	if srvr.httpSrvr.Handler == nil {
		srvr.httpSrvr.Handler = srvr.router
	}

	return &srvr
}

// HTTPServer represents http server implementation
type HTTPServer struct {
	httpSrvr *http.Server
	logger   *log.Logger
	router   *mux.Router
}

// Run will start a server listening on a given port
func (h *HTTPServer) Run(port int) error {
	h.setupServer(port)
	return h.run()
}

func (h *HTTPServer) setupServer(port int) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	h.httpSrvr.Addr = addr
}

func (h *HTTPServer) run() error {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	var err error

	go func() {
		h.logger.Printf("Starting server at: %s", h.httpSrvr.Addr)
		err = h.httpSrvr.ListenAndServe()
	}()

	if err != nil {
		return err
	}

	<-stop

	h.logger.Println("Server shutting down...")
	err = h.Stop()
	if err != nil {
		return err
	}

	h.logger.Println("Server stopped.")

	return nil
}

// Stop attempts to gracefully shutdown the server
func (h *HTTPServer) Stop() error {
	return h.httpSrvr.Shutdown(context.Background())
}

// RegisterServices registers given http Services with
// the server and sets up routes
func (h *HTTPServer) RegisterServices(svcs ...Service) error {
	if svcs == nil {
		return ErrNoServices
	}

	for _, svc := range svcs {
		err := h.RegisterService(svc)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterService registers a given http Service with
// the server and sets up routes
func (h *HTTPServer) RegisterService(svc Service) error {
	endpoints := svc.Endpoints()

	if endpoints == nil {
		return ErrNoEndpointsDefined
	}

	for path, endpoint := range endpoints {
		p := fmt.Sprintf("/%s/%s", svc.Prefix(), path)
		h.router.Handle(p, endpoint.Handler).Methods(endpoint.Methods...)
	}

	return nil
}
