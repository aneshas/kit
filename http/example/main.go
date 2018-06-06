package main

import (
	"log"
	"os"

	"github.com/tonto/kit/http"
	"github.com/tonto/kit/http/adapter"
)

func main() {
	logger := log.New(os.Stdout, "http/example => ", log.Ldate|log.Ltime|log.Lshortfile)

	server := http.NewServer(
		http.WithLogger(logger),
		http.WithAdapters(
			adapter.WithCORS(
				adapter.WithCORSAllowOrigins("*"),
				adapter.WithCORSAllowMethods("PUT", "DELETE"),
				adapter.WithCORSMaxAge(86400),
			),
		),
	)

	reqlog := adapter.WithRequestLogger(logger, true)

	server.MustRegisterServices(
		NewOrderService(reqlog),
		NewCustomerService(reqlog),
	)

	log.Fatal(server.Run(8080))
}
