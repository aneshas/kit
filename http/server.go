// Package server provides common http server functionality
package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

var (
	ErrNoEndpointsDefined = errors.New("service has no endpoints defined")
	ErrNoServices         = errors.New("no services provided")
)

// NewServer creates new http server instance
func NewServer(opts ...ServerOption) *Server {
	srv := Server{
		httpServer: &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			// TODO - Make TLS configurable, maybe even options above
		},
		mux: mux.NewRouter(),
	}

	srv.httpServer.Handler = srv.mux

	for _, o := range opts {
		o(&srv)
	}

	if srv.logger == nil {
		srv.logger = log.New(os.Stdout, "http ", log.Ldate|log.Ltime|log.Llongfile)
	}

	if srv.tls {
		srv.httpServer.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		}
	}

	return &srv
}

// Server represents http server implementation
type Server struct {
	httpServer *http.Server
	logger     *log.Logger
	tls        bool
	certFile   string
	keyFile    string
	mux        *mux.Router
}

// Run will start a server listening on a given port
func (s *Server) Run(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	s.httpServer.Addr = addr

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	var err error

	go func() {
		s.logger.Printf("Starting server at: %s", s.httpServer.Addr)
		if s.tls {
			err = s.runTLS()
		} else {
			err = s.httpServer.ListenAndServe()
		}
	}()

	if err != nil {
		return err
	}

	<-stop

	s.logger.Println("Server shutting down...")
	err = s.Stop()
	if err != nil {
		return err
	}

	s.logger.Println("Server stopped.")

	return nil
}

func (s *Server) runTLS() error {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(srv.ListenAndServe()) }()
	return s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
}

// Stop attempts to gracefully shutdown the server
func (s *Server) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}

// RegisterServices registers given http Services with
// the server and sets up routes
func (s *Server) RegisterServices(svcs ...Service) error {
	if svcs == nil {
		return ErrNoServices
	}

	for _, svc := range svcs {
		err := s.RegisterService(svc)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterService registers a given http Service with
// the server and sets up routes
func (s *Server) RegisterService(svc Service) error {
	endpoints := svc.Endpoints()

	if endpoints == nil {
		return ErrNoEndpointsDefined
	}

	for path, endpoint := range endpoints {
		p := fmt.Sprintf("/%s/%s", svc.Prefix(), path)
		s.mux.Handle(p, endpoint.Handler).Methods(endpoint.Methods...)
	}

	return nil
}
