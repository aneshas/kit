package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	rColor = "\x1b[31;1m"
	gColor = "\x1b[32;1m"
	yColor = "\x1b[33;1m"
	bColor = "\x1b[34;1m"
	wColor = "\x1b[37;1m"
	nColor = "\x1b[0m"
)

// NewServer creates new http server instance
func NewServer(opts ...ServerOption) *Server {
	srv := Server{
		httpServer: &http.Server{
			IdleTimeout: 120 * time.Second,
		},
		stop:         make(chan os.Signal, 1),
		mux:          mux.NewRouter().StrictSlash(true),
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
	}

	if srv.notFoundHandler != nil {
		srv.mux.NotFoundHandler = srv.notFoundHandler
	}

	srv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	srv.httpServer.Handler = srv.mux

	for _, o := range opts {
		o(&srv)
	}

	srv.httpServer.WriteTimeout = srv.writeTimeout
	srv.httpServer.ReadTimeout = srv.readTimeout

	var hf HandlerFunc
	h := srv.httpServer.Handler

	hf = func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(c))
	}

	for _, apt := range srv.adapters {
		hf = apt(hf)
	}

	srv.httpServer.Handler = hf

	if srv.logger == nil {
		srv.logger = log.New(os.Stdout, "kit/http => ", log.Ldate|log.Ltime|log.Llongfile)
	}

	if srv.tlsEnabled() {
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

// Server represents kit http server
type Server struct {
	httpServer      *http.Server
	adapters        []Adapter
	logger          *log.Logger
	certFile        string
	keyFile         string
	mux             *mux.Router
	notFoundHandler http.Handler
	stop            chan os.Signal
	writeTimeout    time.Duration
	readTimeout     time.Duration
}

// Run will start a server listening on a given port
func (s *Server) Run(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	s.httpServer.Addr = addr

	signal.Notify(s.stop, os.Interrupt, os.Kill)

	var err error

	go func() {
		if s.tlsEnabled() {
			s.logger.Printf("%sListening TLS...%s", rColor, nColor)
			s.logger.Println("")
			err = s.runTLS()
			return
		}
		s.logger.Printf("%sServer listening on: %s%s%s", rColor, gColor, s.httpServer.Addr, nColor)
		s.logger.Println("")
		err = s.httpServer.ListenAndServe()
	}()

	<-s.stop

	if err != nil {
		return err
	}

	s.logger.Printf("%sServer shutting down...%s", rColor, nColor)

	e := s.httpServer.Shutdown(context.Background())
	if e != nil {
		return e
	}

	s.logger.Printf("%sServer stopped.%s", rColor, nColor)

	return nil
}

func (s *Server) tlsEnabled() bool {
	return (s.certFile != "" && s.keyFile != "")
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
	s.httpServer.Addr = ":443"
	return s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
}

// Stop attempts to gracefully shutdown the server
func (s *Server) Stop() {
	s.stop <- os.Interrupt
}

// RegisterServices registers given http Services with
// the server and sets up routes
func (s *Server) RegisterServices(svcs ...Service) error {
	if svcs == nil {
		return fmt.Errorf("no services provided")
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
		return fmt.Errorf("service has no endpoints defined")
	}

	st := reflect.ValueOf(svc).Elem().Type()
	s.logger.Printf(
		"%sREGISTERING %s%s:%s%s %s%s",
		wColor,
		bColor, st.PkgPath(),
		yColor, st.Name(),
		wColor,
		nColor,
	)

	for path, endpoint := range endpoints {
		s.printRouteInfo(svc, path, endpoint)
		hfunc := endpoint.Handler

		route := s.mux.HandleFunc(
			s.getPath(path, svc.Prefix()),
			func(w http.ResponseWriter, r *http.Request) {
				// TODO - Provide a sensible default context
				// eg. timeouts, values ???
				hfunc(r.Context(), w, r)
			},
		)

		if endpoint.Methods != nil {
			route.Methods(endpoint.Methods...)
		}
	}

	s.logger.Println("")

	return nil
}

func (s *Server) printRouteInfo(svc Service, path string, ep *Endpoint) {
	s.logger.Printf(
		"%s%s %s/%s %s",
		yColor, strings.Join(ep.Methods, ","),
		gColor, svc.Prefix()+path,
		nColor,
	)
}

func (s *Server) getPath(path string, prefix string) string {
	if path == "/" {
		path = ""
	}
	return fmt.Sprintf("/%s%s", prefix, strings.TrimRight(path, "/"))
}
